package rntocase

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameFilesWithDiscovery_Collision(t *testing.T) {
	tempDir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_one.txt"), []byte("test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_two.txt"), []byte("test"), 0644))

	renameFunc := func(s string) (string, error) {
		return "renamed", nil // They will both map to `renamed.txt` in the same tempDir namespace
	}

	err := RenameFilesWithDiscovery([]string{tempDir}, true, nil, nil, renameFunc, false, false, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "collision")
}

func TestRenameFilesWithDiscovery_DryRun(t *testing.T) {
	tempDir := t.TempDir()

	f1 := filepath.Join(tempDir, "file_one.txt")
	require.NoError(t, os.WriteFile(f1, []byte("test"), 0644))

	subDir := filepath.Join(tempDir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))
	f2 := filepath.Join(subDir, "file_two.txt")
	require.NoError(t, os.WriteFile(f2, []byte("test"), 0644))

	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	err := RenameFilesWithDiscovery([]string{tempDir}, true, nil, nil, renameFunc, true, false, false)
	require.NoError(t, err)

	// Ensure files are not changed
	_, err1 := os.Stat(f1)
	require.NoError(t, err1)

	_, err2 := os.Stat(f2)
	require.NoError(t, err2)
}
