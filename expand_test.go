package rntocase

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandFiles(t *testing.T) {
	tempDir := t.TempDir()

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

	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file2.jpg"), []byte("test"), 0644))

	subDir := filepath.Join(tempDir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "file4.jpg"), []byte("test"), 0644))

	gitDir := filepath.Join(subDir, ".git")
	require.NoError(t, os.Mkdir(gitDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(gitDir, "config"), []byte("test"), 0644))

	require.NoError(t, os.Symlink("sub", filepath.Join(tempDir, "sym_dir")))
	require.NoError(t, os.Symlink("file1.txt", filepath.Join(tempDir, "sym_file")))

	t.Run("NonRecursive", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{tempDir}, files)
	})

	t.Run("Recursive No Filters", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, true, nil, nil)
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
		files, err := ExpandFiles([]string{tempDir}, true, []string{"*.jpg"}, nil)
		require.NoError(t, err)

		expected := []string{
			filepath.Join(tempDir, "file2.jpg"),
			filepath.Join(tempDir, "sub", "file4.jpg"),
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Recursive With Exclude Git", func(t *testing.T) {
		files, err := ExpandFiles([]string{tempDir}, true, nil, []string{".git/**"})
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
		files, err := ExpandFiles([]string{tempDir, subDir}, true, nil, []string{".git/**"})
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
		files, err := ExpandFiles([]string{f}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{f}, files)
	})

	t.Run("NonRecursive Symlink File", func(t *testing.T) {
		f := filepath.Join(tempDir, "sym_file")
		files, err := ExpandFiles([]string{f}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{f}, files)
	})

	t.Run("Space in path", func(t *testing.T) {
		spaceDir := filepath.Join(tempDir, "space dir")
		require.NoError(t, os.Mkdir(spaceDir, 0755))
		f := filepath.Join(spaceDir, "space file.txt")
		require.NoError(t, os.WriteFile(f, []byte("test"), 0644))

		files, err := ExpandFiles([]string{spaceDir}, true, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{f}, files)
	})
}

func TestExpandFiles_ExplicitFileRoot(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "file1.txt")
	require.NoError(t, os.WriteFile(file1, []byte("test"), 0644))

	file2 := filepath.Join(tempDir, "file2.jpg")
	require.NoError(t, os.WriteFile(file2, []byte("test"), 0644))

	t.Run("Recursive With Include JPG", func(t *testing.T) {
		files, err := ExpandFiles([]string{file1, file2}, true, []string{"*.jpg"}, nil)
		require.NoError(t, err)

		expected := []string{
			file2,
		}

		assert.Equal(t, expected, files)
	})

	t.Run("Recursive With Exclude TXT", func(t *testing.T) {
		files, err := ExpandFiles([]string{file1, file2}, true, nil, []string{"*.txt"})
		require.NoError(t, err)

		expected := []string{
			file2,
		}

		assert.Equal(t, expected, files)
	})

	t.Run("NonRecursive Preserves Missing", func(t *testing.T) {
		files, err := ExpandFiles([]string{"missing_file"}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{"missing_file"}, files)
	})
}

func TestExpandFiles_SymlinkRoots(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "file1.txt")
	require.NoError(t, os.WriteFile(file1, []byte("test"), 0644))

	symFile := filepath.Join(tempDir, "sym_file.txt")
	require.NoError(t, os.Symlink(file1, symFile))

	subDir := filepath.Join(tempDir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))

	symDir := filepath.Join(tempDir, "sym_dir")
	require.NoError(t, os.Symlink(subDir, symDir))

	t.Run("Recursive Skips Explicit File Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symFile}, true, nil, nil)
		require.NoError(t, err)
		assert.Empty(t, files) // should skip completely
	})

	t.Run("Recursive Skips Explicit Dir Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symDir}, true, nil, nil)
		require.NoError(t, err)
		assert.Empty(t, files) // should skip completely
	})

	t.Run("NonRecursive Preserves File Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symFile}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{symFile}, files)
	})

	t.Run("NonRecursive Preserves Dir Symlink", func(t *testing.T) {
		files, err := ExpandFiles([]string{symDir}, false, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, []string{symDir}, files)
	})
}
