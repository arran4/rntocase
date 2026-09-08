package skill

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTarball(t *testing.T, files map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.tar.gz")
	assert.NoError(t, err)

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	assert.NoError(t, tw.Close())
	assert.NoError(t, gw.Close())
	assert.NoError(t, f.Close())

	return f.Name()
}

func TestExtractTarGz_PathTraversal(t *testing.T) {
	destDir := t.TempDir()

	files := map[string]string{
		"repo-sha/valid.txt":      "valid content",
		"repo-sha/../../evil.txt": "evil content",
	}

	tarPath := createTestTarball(t, files)
	defer func() { _ = os.Remove(tarPath) }()

	err := ExtractTarGz(tarPath, destDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal detected")
}

func TestExtractTarGz_Success(t *testing.T) {
	destDir := t.TempDir()

	files := map[string]string{
		"repo-sha/SKILL.md":      "# My Skill\n",
		"repo-sha/lib/helper.py": "print('hello')",
	}

	tarPath := createTestTarball(t, files)
	defer func() { _ = os.Remove(tarPath) }()

	err := ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "# My Skill\n", string(content))
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
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "stale.txt"), []byte("stale"), 0644)
	assert.NoError(t, err)

	err = ReplaceSafely(destDir, func(stagingDir string) error {
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
