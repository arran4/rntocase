package rntocase

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	_ "embed"
	"github.com/arran4/rntocase/internal/fsys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
)

//go:embed testdata/expand/expand.txtar
var expandTxtar []byte

func getExpected(t *testing.T, ar *txtar.Archive, name string) []string {
	for _, f := range ar.Files {
		if f.Name == name {
			var res []string
			err := json.Unmarshal(f.Data, &res)
			require.NoError(t, err)
			return res
		}
	}
	t.Fatalf("expected file %s not found in txtar", name)
	return nil
}

func getMockFS(t *testing.T, ar *txtar.Archive) *fsys.MockFS {
	fs := fsys.NewMockFS()
	for _, f := range ar.Files {
		if strings.HasPrefix(f.Name, "fs/") {
			name := "/" + strings.TrimPrefix(f.Name, "fs/")
			if strings.HasSuffix(name, "/") {
				if strings.HasPrefix(string(f.Data), "SYMLINK: ") {
					fs.AddSymlink(name[:len(name)-1], strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK: ")))
				} else {
					fs.AddDir(name)
				}
			} else {
				if strings.HasPrefix(string(f.Data), "SYMLINK: ") {
					fs.AddSymlink(name, strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK: ")))
				} else {
					fs.AddFile(name)
				}
			}
		}
	}
	return fs
}

func TestExpandFiles(t *testing.T) {
	ar := txtar.Parse(expandTxtar)
	fs := getMockFS(t, ar)
	tempDir := "/testdir"
	fs.AddDir("/space dir")
	fs.AddFile("/space dir/space file.txt")

	// Create test structure
	// tempDir/
	//   file1.txt
	//   file2.jpg
	//   sub/
	//     file3.txt
	//     file4.jpg
	//     .git/
	//       config
	//   sym_dir -> sub
	//   sym_file -> file1.txt

	subDir := filepath.Join(tempDir, "sub")

	t.Run("NonRecursive", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{tempDir}, files)
	})

	t.Run("Recursive No Filters", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, true, nil, nil, fs)
		require.NoError(t, err)

		expected := []string{
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "file2.jpg"),
			filepath.Join(tempDir, "sub", ".git", "config"),
			filepath.Join(tempDir, "sub", "file3.txt"),
			filepath.Join(tempDir, "sub", "file4.jpg"),
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Recursive With Include JPG", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, true, []string{"*.jpg"}, nil, fs)
		require.NoError(t, err)

		expected := []string{
			filepath.Join(tempDir, "file2.jpg"),
			filepath.Join(tempDir, "sub", "file4.jpg"),
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Recursive With Exclude Git", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, true, nil, []string{".git/**"}, fs)
		require.NoError(t, err)

		expected := []string{
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "file2.jpg"),
			filepath.Join(tempDir, "sub", "file3.txt"),
			filepath.Join(tempDir, "sub", "file4.jpg"),
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Overlapping Roots Deduplication", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir, subDir}, true, nil, []string{".git/**"}, fs)
		require.NoError(t, err)

		expected := []string{
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "file2.jpg"),
			filepath.Join(tempDir, "sub", "file3.txt"),
			filepath.Join(tempDir, "sub", "file4.jpg"),
		}

		assert.Equal(t, expected, files)
	})

	t.Run("NonRecursive File Exists", func(t *testing.T) {
		f := filepath.Join(tempDir, "file1.txt")
		files, err := ExpandFiles([]string{f}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{f}, files)
	})

	t.Run("NonRecursive Symlink File", func(t *testing.T) {
		f := filepath.Join(tempDir, "sym_file")
		files, err := ExpandFiles([]string{f}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{f}, files)
	})

	t.Run("Space in path", func(t *testing.T) {
		spaceDir := "/space dir"
		fs.AddDir("/space dir")
		fs.AddFile("/space dir/space file.txt")

		files, err := ExpandFiles([]string{spaceDir}, true, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{"/space dir/space file.txt"}, files)
	})
}

func TestExpandFiles_ExplicitFileRoot(t *testing.T) {
	fs := fsys.NewMockFS()
	tempDir := "/testdir"

	file1 := filepath.Join(tempDir, "file1.txt")
	fs.AddFile(file1)

	file2 := filepath.Join(tempDir, "file2.jpg")
	fs.AddFile(file2)

	t.Run("Recursive With Include JPG", func(t *testing.T) {
		files, err := ExpandFiles([]string{file1, file2}, true, []string{"*.jpg"}, nil, fs)
		require.NoError(t, err)

		expected := []string{
			file2,
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Recursive With Exclude TXT", func(t *testing.T) {
		files, err := ExpandFiles([]string{file1, file2}, true, nil, []string{"*.txt"}, fs)
		require.NoError(t, err)

		expected := []string{
			file2,
		}

		assert.Equal(t, expected, files)
	})

	t.Run("NonRecursive Preserves Missing", func(t *testing.T) {
		files, err := ExpandFiles([]string{"missing_file"}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{"missing_file"}, files)
	})
}

func TestExpandFiles_SymlinkRoots(t *testing.T) {
	fs := fsys.NewMockFS()
	tempDir := "/testdir"

	file1 := filepath.Join(tempDir, "file1.txt")
	fs.AddFile(file1)

	symFile := filepath.Join(tempDir, "sym_file.txt")
	fs.AddSymlink(symFile, "/testdir/file1.txt")

	subDir := filepath.Join(tempDir, "sub")
	fs.AddDir(subDir)

	symDir := filepath.Join(tempDir, "sym_dir")
	fs.AddSymlink(symDir, "/testdir/sub")

	t.Run("Recursive Skips Explicit File Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symFile}, true, nil, nil, fs)
		require.NoError(t, err)
		assert.Empty(t, files) // should skip completely
	})

	t.Run("Recursive Skips Explicit Dir Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symDir}, true, nil, nil, fs)
		require.NoError(t, err)
		assert.Empty(t, files) // should skip completely
	})

	t.Run("NonRecursive Preserves File Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symFile}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{symFile}, files)
	})

	t.Run("NonRecursive Preserves Dir Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symDir}, false, nil, nil, fs)
		require.NoError(t, err)
		assert.Equal(t, []string{symDir}, files)
	})
}
