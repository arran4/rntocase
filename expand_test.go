package rntocase

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"embed"
	"github.com/arran4/rntocase/internal/fstestutil"
	"github.com/arran4/rntocase/internal/fsys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
)

//go:embed testdata/expand/cases/*.txtar
var expandTxtar embed.FS

type txtarOptions struct {
	Roots     []string `json:"roots"`
	Recursive bool     `json:"recursive"`
	Includes  []string `json:"includes"`
	Excludes  []string `json:"excludes"`
}

func TestExpandFiles_Txtar(t *testing.T) {
	entries, err := fs.Glob(expandTxtar, "testdata/expand/cases/*.txtar")
	require.NoError(t, err)
	require.NotEmpty(t, entries, "require at least one embedded scenario")

	for _, fixture := range entries {
		fixture := fixture
		t.Run(strings.TrimSuffix(filepath.Base(fixture), ".txtar"), func(t *testing.T) {
			raw, err := fs.ReadFile(expandTxtar, fixture)
			require.NoError(t, err)

			ar := txtar.Parse(raw)
			mockFs := fstestutil.ArchiveToMockFS(ar)

			var opts txtarOptions
			var expected []string

			for _, f := range ar.Files {
				switch f.Name {
				case "options.json":
					err = json.Unmarshal(f.Data, &opts)
					require.NoError(t, err)
				case "expected.json":
					err = json.Unmarshal(f.Data, &expected)
					require.NoError(t, err)
				}
			}

			if expected == nil {
				t.Fatalf("missing expected.json in %s", fixture)
			}
			// require options.json to have been parsed, but we don't have a direct check except maybe we can verify it was in the archive
			var hasOptions bool
			for _, f := range ar.Files {
				if f.Name == "options.json" {
					hasOptions = true
				}
			}
			if !hasOptions {
				t.Fatalf("missing options.json in %s", fixture)
			}

			files, err := ExpandFiles(opts.Roots, opts.Recursive, opts.Includes, opts.Excludes, mockFs)
			require.NoError(t, err)

			if len(expected) == 0 && len(files) == 0 {
				// both empty, ok
			} else {
				assert.Equal(t, expected, files)
			}
		})
	}
}

type errorFS struct {
	fsys.MockFS
}

func (e *errorFS) Lstat(name string) (os.FileInfo, error) {
	if name == "/testdir" {
		return nil, os.ErrPermission
	}
	return e.MockFS.Lstat(name)
}

func TestExpandFiles_InfoError(t *testing.T) {
	fs := fsys.NewMockFS()
	fs.AddFile("/testdir/error_trigger.txt")
	fs.AddFile("/testdir/normal.txt")

	errFs := &errorFS{*fs}

	_, err := ExpandFiles([]string{"/testdir"}, true, nil, nil, errFs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission")
}

func TestExpandFiles_FailClosedInvalidInjection(t *testing.T) {
	_, err := ExpandFiles([]string{"foo.txt"}, false, nil, nil, struct{}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported filesystem injected")

	var nilMock *fsys.MockFS = nil
	_, err = ExpandFiles([]string{"foo.txt"}, false, nil, nil, nilMock)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported filesystem injected")
}

func TestRenameFilesWithDiscovery_Memory(t *testing.T) {
	ar := txtar.Parse([]byte("-- fs/testdir/file1.txt --\n-- fs/testdir/file2.jpg --\n-- fs/testdir/sub/file3.txt --\n-- fs/testdir/sub/file4.jpg --\n-- fs/testdir/sub/.git/config --"))
	fs := fstestutil.ArchiveToMockFS(ar)

	// We want to verify end to end against mockfs
	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	err := RenameFilesWithDiscovery([]string{"/testdir"}, true, nil, nil, renameFunc, false, false, false, fs)
	require.NoError(t, err)

	expectedFiles := []string{
		"/testdir/file1_renamed.txt",
		"/testdir/file2_renamed.jpg",
		"/testdir/sub/.git/config_renamed", // actually config has no extension so config_renamed
		"/testdir/sub/file3_renamed.txt",
		"/testdir/sub/file4_renamed.jpg",
	}

	for _, f := range expectedFiles {
		_, err := fs.Stat(f)
		assert.NoError(t, err, "Expected file %s to exist in MockFS after rename", f)
	}
}

func TestRenameFilesWithDiscovery_Memory_Collision(t *testing.T) {
	fs := fsys.NewMockFS()
	fs.AddFile("/testdir/file1.txt")
	fs.AddFile("/testdir/file2.txt")

	renameFunc := func(s string) (string, error) {
		return "renamed", nil
	}

	err := RenameFilesWithDiscovery([]string{"/testdir"}, true, nil, nil, renameFunc, false, false, false, fs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "collision")

	// Ensure old files are preserved
	_, err1 := fs.Stat("/testdir/file1.txt")
	assert.NoError(t, err1)
	_, err2 := fs.Stat("/testdir/file2.txt")
	assert.NoError(t, err2)
}

func TestArchiveToMockFS_MetadataOnly(t *testing.T) {
	ar := txtar.Parse([]byte("-- options.json --\n{\"key\":\"value\"}\n-- expected.json --\n[\"foo\"]"))
	fs := fstestutil.ArchiveToMockFS(ar)

	// Since there is no fs/ files, MockFS should have no files other than the root / (which it might create lazily or stay empty)
	if len(fs.Files) > 0 {
		t.Fatalf("Expected empty mockFS, got %d files", len(fs.Files))
	}
}
