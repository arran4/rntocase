package skill

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(path, 0755))
}

func mustWriteFile(t *testing.T, filename, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filename, []byte(content), 0644))
}

func mustSymlink(t *testing.T, oldname, newname string) {
	t.Helper()
	require.NoError(t, os.Symlink(oldname, newname))
}

func createTestTarball(t *testing.T, files map[string]string) string {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), "test.tar.gz"))
	require.NoError(t, err)

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	// Sort keys to ensure deterministic ordering (e.g. directories before symlinks).
	// String sort will naturally place 'repo-sha/sub/' before 'repo-sha/sub/link'.
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		content := files[name]
		var typeflag byte = tar.TypeReg
		linkname := ""
		if strings.HasSuffix(name, "/") {
			typeflag = tar.TypeDir
		} else if strings.HasPrefix(content, "symlink:") {
			typeflag = tar.TypeSymlink
			linkname = strings.TrimPrefix(content, "symlink:")
			content = ""
		}
		hdr := &tar.Header{
			Name:     name,
			Mode:     0600,
			Size:     int64(len(content)),
			Typeflag: typeflag,
			Linkname: linkname,
		}
		require.NoError(t, tw.WriteHeader(hdr))
		if content != "" {
			_, err := tw.Write([]byte(content))
			require.NoError(t, err)
		}
	}

	require.NoError(t, tw.Close())
	require.NoError(t, gw.Close())
	require.NoError(t, f.Close())

	return f.Name()
}

func mockGitHubAPI(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(handler)
	originalAPIURL := GitHubAPIURL
	GitHubAPIURL = ts.URL
	t.Cleanup(func() {
		ts.Close()
		GitHubAPIURL = originalAPIURL
	})
	return ts
}

func mockHTTPClientTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()
	originalTimeout := HTTPClient.Timeout
	HTTPClient.Timeout = timeout
	t.Cleanup(func() {
		HTTPClient.Timeout = originalTimeout
	})
}

func TestNetwork_Timeout(t *testing.T) {
	// Setup a slow mock server that hangs for 100ms
	mockGitHubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	mockHTTPClientTimeout(t, 10*time.Millisecond) // Shorter than the server sleep

	// Test DownloadGitHubRepository timeout
	_, _, err := DownloadGitHubRepository("dummy/repo", "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded, "expected deadline exceeded error")

	// Test CheckUpdate timeout
	meta := &Metadata{OwnerRepo: "dummy/repo"}
	_, _, err = CheckUpdate(meta)
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded, "expected deadline exceeded error")
}

func TestNetwork_Non200(t *testing.T) {
	// Setup a mock server that returns 500
	mockGitHubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	// Test DownloadGitHubRepository non-200
	_, _, err := DownloadGitHubRepository("dummy/repo", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")

	// Test CheckUpdate non-200
	meta := &Metadata{OwnerRepo: "dummy/repo"}
	_, _, err = CheckUpdate(meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

func TestExtractTarGz_PathTraversal(t *testing.T) {
	destDir := t.TempDir()

	files := map[string]string{
		"repo-sha/valid.txt":      "valid content",
		"repo-sha/../../evil.txt": "evil content",
	}

	tarPath := createTestTarball(t, files)

	err := ExtractTarGz(tarPath, destDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal detected")
}

func TestExtractTarGz_PathTraversalExactRoot(t *testing.T) {
	destDir := t.TempDir()

	// Test extraction that directly creates a symlink within the destination folder
	// whose resolved path is exactly the destination root folder, which would previously
	// fail the `HasPrefix` test as there is no trailing slash.
	files := map[string]string{
		"repo-sha/sub/":     "",
		"repo-sha/sub/link": "symlink:../",
	}

	tarPath := createTestTarball(t, files)

	// Should extract successfully since the symlink points safely within destDir.
	err := ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	linkTarget, err := os.Readlink(filepath.Join(destDir, "sub", "link"))
	assert.NoError(t, err)
	assert.Equal(t, "../", linkTarget)
}

func TestValidateSymlinkTarget_WindowsCases(t *testing.T) {
	destDir := "/opt/dest"
	targetPath := "/opt/dest/sub/link"

	tests := []struct {
		name          string
		symlinkTarget string
		wantErr       bool
		errContains   string
	}{
		{"POSIX absolute", "/etc/passwd", true, "symlinks"},
		{"Windows volume absolute (POSIX context)", "C:\\Windows", false, ""},             // Valid literal filename on POSIX
		{"Windows volume qualified relative (POSIX context)", "C:..\\outside", false, ""}, // Valid literal filename on POSIX
		{"Windows rooted path (POSIX context)", "\\outside", false, ""},                   // Valid literal filename on POSIX
		{"POSIX rooted path equivalent", "/outside", true, "symlinks"},                    // Will hit IsAbs first on POSIX usually
		{"Valid POSIX relative", "target.txt", false, ""},
		{"Valid exact root", "../", false, ""},
		{"Valid Windows exact root", "..\\", false, ""},
		{"Escaping POSIX relative", "../../outside", true, "symlink points outside"},
		{"Escaping Windows relative (POSIX context)", "..\\..\\outside", false, ""}, // On POSIX, backslashes are literal filenames, so this is valid.
	}

	if runtime.GOOS == "windows" {
		for i, tt := range tests {
			switch tt.name {
			case "Escaping Windows relative (POSIX context)":
				tests[i].name = "Escaping Windows relative"
				tests[i].wantErr = true
				tests[i].errContains = "symlink points outside"
			case "Windows volume absolute (POSIX context)":
				tests[i].name = "Windows volume absolute"
				tests[i].wantErr = true
				tests[i].errContains = "symlinks"
			case "Windows volume qualified relative (POSIX context)":
				tests[i].name = "Windows volume qualified relative"
				tests[i].wantErr = true
				tests[i].errContains = "volume-qualified symlinks"
			case "Windows rooted path (POSIX context)":
				tests[i].name = "Windows rooted path"
				tests[i].wantErr = true
				tests[i].errContains = "rooted symlinks"
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSymlinkTarget(tt.symlinkTarget, targetPath, destDir)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractTarGz_Success(t *testing.T) {
	destDir := t.TempDir()

	files := map[string]string{
		"repo-sha/SKILL.md":      "---\nname: my-skill\ndescription: test\n---\n# My Skill\n",
		"repo-sha/lib/helper.py": "print('hello')",
	}

	tarPath := createTestTarball(t, files)

	err := ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "---\nname: my-skill\ndescription: test\n---\n# My Skill\n", string(content))
}

func TestReplaceSafelyWithFS_Success(t *testing.T) {
	fakeFS := NewFakeFSOps()
	destDir := "/opt/skills/my-skill"

	// Initially does not exist
	err := replaceSafelyWithFS(destDir, func(stagingDir string) error {
		fakeFS.Files[filepath.Join(stagingDir, "SKILL.md")] = true
		return nil
	}, fakeFS)
	assert.NoError(t, err)

	assert.True(t, fakeFS.Dirs[destDir])
	assert.True(t, fakeFS.Files[filepath.Join(destDir, "SKILL.md")])

	// Replace existing
	err = replaceSafelyWithFS(destDir, func(stagingDir string) error {
		fakeFS.Files[filepath.Join(stagingDir, "NEW_FILE.md")] = true
		return nil
	}, fakeFS)
	assert.NoError(t, err)

	assert.True(t, fakeFS.Dirs[destDir])
	assert.False(t, fakeFS.Files[filepath.Join(destDir, "SKILL.md")]) // Old file gone
	assert.True(t, fakeFS.Files[filepath.Join(destDir, "NEW_FILE.md")])
}

func TestReplaceSafelyWithFS_RollbackOnFailure(t *testing.T) {
	fakeFS := NewFakeFSOps()
	destDir := "/opt/skills/my-skill"
	fakeFS.Dirs[destDir] = true
	fakeFS.Files[filepath.Join(destDir, "SKILL.md")] = true

	// Simulate a failure during populate (e.g., validation failed)
	err := replaceSafelyWithFS(destDir, func(stagingDir string) error {
		return os.ErrPermission // Some simulated error
	}, fakeFS)
	assert.Error(t, err)

	// Original destination should remain completely untouched
	assert.True(t, fakeFS.Dirs[destDir])
	assert.True(t, fakeFS.Files[filepath.Join(destDir, "SKILL.md")])
}

func TestReplaceSafelyWithFS_RemovesStaleFiles(t *testing.T) {
	fakeFS := NewFakeFSOps()
	destDir := "/opt/skills/my-skill"
	fakeFS.Dirs[destDir] = true
	fakeFS.Files[filepath.Join(destDir, "SKILL.md")] = true
	fakeFS.Files[filepath.Join(destDir, "stale.py")] = true

	// Replace existing but new version doesn't have stale.py
	err := replaceSafelyWithFS(destDir, func(stagingDir string) error {
		fakeFS.Files[filepath.Join(stagingDir, "SKILL.md")] = true
		return nil
	}, fakeFS)
	assert.NoError(t, err)

	// Verify old files are gone
	assert.False(t, fakeFS.Files[filepath.Join(destDir, "stale.py")])
	assert.True(t, fakeFS.Files[filepath.Join(destDir, "SKILL.md")])
}

func TestReplaceSafelyWithFS_CommitRenameFailureRollback(t *testing.T) {
	fakeFS := NewFakeFSOps()
	destDir := "/opt/skills/my-skill"
	fakeFS.Dirs[destDir] = true
	fakeFS.Files[filepath.Join(destDir, "SKILL.md")] = true

	var stagingDirName string
	err := replaceSafelyWithFS(destDir, func(stagingDir string) error {
		stagingDirName = stagingDir
		fakeFS.Files[filepath.Join(stagingDir, "new-v2")] = true

		// Inject failure for commit phase: moving stagingDir -> destDir fails
		fakeFS.RenameFailFS[stagingDir+"->"+destDir] = os.ErrPermission
		return nil
	}, fakeFS)

	// Expect the commit to fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit new installation")
	// Ensure it didn't throw a rollback failed error
	assert.NotContains(t, err.Error(), "rollback failed")

	// Verify rollback succeeded: original files restored
	assert.True(t, fakeFS.Dirs[destDir])
	assert.True(t, fakeFS.Files[filepath.Join(destDir, "SKILL.md")])
	assert.False(t, fakeFS.Files[filepath.Join(destDir, "new-v2")])

	// Staging dir is cleaned up via defer
	assert.False(t, fakeFS.Dirs[stagingDirName])
}

// Keep one OS integration test to ensure osRename / ReplaceSafely actually works on real disk
func TestReplaceSafely_Integration(t *testing.T) {
	parentDir := t.TempDir()
	destDir := filepath.Join(parentDir, "my-skill")
	mustMkdirAll(t, destDir)
	mustWriteFile(t, filepath.Join(destDir, "stale.txt"), "stale")

	err := ReplaceSafely(destDir, func(stagingDir string) error {
		return os.WriteFile(filepath.Join(stagingDir, "SKILL.md"), []byte("v2"), 0644)
	})
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "v2", string(content))

	_, err = os.Stat(filepath.Join(destDir, "stale.txt"))
	assert.True(t, os.IsNotExist(err))
}

// Fixing the double failure test to use the fake properly
func TestReplaceSafelyWithFS_DoubleFailure(t *testing.T) {
	fakeFS := NewFakeFSOps()
	destDir := "/opt/skills/my-skill"
	fakeFS.Dirs[destDir] = true
	fakeFS.Files[filepath.Join(destDir, "SKILL.md")] = true

	// Custom Rename in FakeFSOps to fail both commit and rollback
	originalRename := fakeFS.RenameFailFS
	defer func() { fakeFS.RenameFailFS = originalRename }()

	// For the fake FS we just need to ensure Rename returns error when target is destDir
	err := replaceSafelyWithFS(destDir, func(stagingDir string) error {
		// Just fail any rename that targets destDir
		fakeFS.RenameFailFS[stagingDir+"->"+destDir] = os.ErrPermission
		// Note: The backup dir name is dynamically generated. In fakeFS, MkdirTemp yields predictable names like `/opt/skills/.backup-skill-X`
		// We'll just hardcode a catch-all in the Rename logic if we need to, or just find the backup dir.
		return nil
	}, &failAllDestRenameFS{FakeFSOps: fakeFS, destDir: destDir})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit new installation")
	assert.Contains(t, err.Error(), "and rollback failed")
}

type failAllDestRenameFS struct {
	*FakeFSOps
	destDir string
}

func (f *failAllDestRenameFS) Rename(oldpath, newpath string) error {
	if newpath == f.destDir {
		return os.ErrPermission
	}
	return f.FakeFSOps.Rename(oldpath, newpath)
}

func TestDownloadGitHubRepository_CustomRef(t *testing.T) {
	var capturedURL string
	mockGitHubAPI(t, func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.Path
		if strings.HasSuffix(r.URL.Path, "tarball/sha-456") {
			w.WriteHeader(http.StatusOK)
			// Need a valid tarball to avoid gzip reader errors, or just let it fail at extract
			// We only care about the URL for this test, so writing some junk and ignoring extraction error is ok.
			return
		}

		commit := GitHubCommit{Sha: "sha-456"}
		_ = json.NewEncoder(w).Encode(commit)
	})

	// Even if it fails creating/writing tarball due to empty response body, we just check if it queried the right path
	_, _, err := DownloadGitHubRepository("dummy/repo", "v1.0")
	_ = err
	assert.Contains(t, capturedURL, "tarball/sha-456")
}

func TestExtractTarGz_MissingSKILLmd(t *testing.T) {
	destDir := t.TempDir()

	files := map[string]string{
		"repo-sha/lib/helper.py": "print('hello')",
	}

	tarPath := createTestTarball(t, files)

	err := ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(destDir, "SKILL.md"))
	assert.True(t, os.IsNotExist(err))
}

func TestClassifySource(t *testing.T) {
	tempDir := t.TempDir()

	// Change dir to temp root for relative path testing
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer func() {
		err := os.Chdir(oldWd)
		require.NoError(t, err)
	}()
	err = os.MkdirAll("skills/example", 0755)
	require.NoError(t, err)

	tests := []struct {
		name          string
		source        string
		wantLocal     bool
		wantOfficial  bool
		wantOwnerRepo string
		wantErr       bool
	}{
		{"official", "official", false, true, "", false},
		{"rntocase", "rntocase", false, true, "", false},
		{"owner/repo", "owner/repo", false, false, "owner/repo", false},
		{"invalid remote", "owner-repo", false, false, "", true},
		{"existing relative dir", "skills/example", true, false, "", false},
		{"explicit missing relative", "./missing-dir", false, false, "", true},
		{"explicit missing parent", "../missing-dir", false, false, "", true},
		{"explicit missing absolute", "/missing-dir-123", false, false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isLocal, isOfficial, ownerRepo, err := ClassifySource(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassifySource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantLocal, isLocal)
			assert.Equal(t, tt.wantOfficial, isOfficial)
			assert.Equal(t, tt.wantOwnerRepo, ownerRepo)
		})
	}
}

func TestCopyLocalDirectory_Symlinks(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(srcDir, "sub"))
	mustWriteFile(t, filepath.Join(srcDir, "sub", "target.txt"), "hello")

	// Test case 1: Valid symlink contained within the copied skill
	mustSymlink(t, "target.txt", filepath.Join(srcDir, "sub", "valid_link"))

	// Test case 2: Relative symlink escaping the destination
	mustSymlink(t, "../../outside", filepath.Join(srcDir, "sub", "escaping_link"))

	err := CopyLocalDirectory(srcDir, destDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symlink points outside destination directory")
}

func TestCopyLocalDirectory_AbsoluteSymlinks(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(srcDir, "sub"))

	// Test case 3: Absolute outside symlink
	mustSymlink(t, "/etc/passwd", filepath.Join(srcDir, "sub", "abs_link"))

	err := CopyLocalDirectory(srcDir, destDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symlinks")
}

func TestCopyLocalDirectory_ValidSymlinkOnly(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(srcDir, "sub"))
	mustWriteFile(t, filepath.Join(srcDir, "sub", "target.txt"), "hello")

	// Valid symlink contained within the copied skill
	mustSymlink(t, "target.txt", filepath.Join(srcDir, "sub", "valid_link"))

	err := CopyLocalDirectory(srcDir, destDir)
	assert.NoError(t, err)

	// Verify the symlink was created correctly
	linkTarget, err := os.Readlink(filepath.Join(destDir, "sub", "valid_link"))
	assert.NoError(t, err)
	assert.Equal(t, "target.txt", linkTarget)
}

func TestCopyLocalDirectory_ExactRootSymlink(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	mustMkdirAll(t, filepath.Join(srcDir, "sub"))

	// Test case 4: Valid exact-root symlink
	mustSymlink(t, "../", filepath.Join(srcDir, "sub", "exact_root_link"))

	err := CopyLocalDirectory(srcDir, destDir)
	assert.NoError(t, err)

	// Verify the symlink was created correctly
	linkTarget, err := os.Readlink(filepath.Join(destDir, "sub", "exact_root_link"))
	assert.NoError(t, err)
	assert.Equal(t, "../", linkTarget)
}
