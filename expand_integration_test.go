package rntocase

import (
	"bytes"
	"encoding/json"
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

func TestRenameFilesWithDiscovery_MissingRootRecursive(t *testing.T) {
	tempDir := t.TempDir()

	validFile := filepath.Join(tempDir, "file_one.txt")
	require.NoError(t, os.WriteFile(validFile, []byte("test"), 0644))

	missingRoot := filepath.Join(tempDir, "does_not_exist")

	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	err := RenameFilesWithDiscovery([]string{validFile, missingRoot}, true, nil, nil, renameFunc, false, false, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "recursive root not accessible")

	// Ensure valid file was NOT mutated because preflight fails early
	_, err1 := os.Stat(validFile)
	require.NoError(t, err1)

	renamedFile := filepath.Join(tempDir, "file_one_renamed.txt")
	_, err2 := os.Stat(renamedFile)
	require.Error(t, err2)
	assert.True(t, os.IsNotExist(err2))
}

func TestRenameFilesWithDiscovery_JSONOutput(t *testing.T) {
	tempDir := t.TempDir()

	f1 := filepath.Join(tempDir, "file_one.txt")
	require.NoError(t, os.WriteFile(f1, []byte("test"), 0644))

	f2 := filepath.Join(tempDir, "file_two.txt")
	require.NoError(t, os.WriteFile(f2, []byte("test"), 0644))

	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	// Capture stdout safely
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Stdout = oldStdout
		_ = w.Close()
		_ = r.Close()
	})
	os.Stdout = w

	outCh := make(chan []byte, 1)
	readErrCh := make(chan error, 1)

	go func() {
		var buf bytes.Buffer
		_, readErr := buf.ReadFrom(r)
		readErrCh <- readErr
		outCh <- buf.Bytes()
	}()

	renameErr := RenameFilesWithDiscovery([]string{tempDir}, true, nil, nil, renameFunc, true, false, true) // json=true, dryRun=true

	os.Stdout = oldStdout
	require.NoError(t, w.Close())

	require.NoError(t, <-readErrCh)
	capturedBytes := <-outCh

	require.NoError(t, r.Close())
	require.NoError(t, renameErr)

	var result RenameResult
	err = json.Unmarshal(capturedBytes, &result)
	require.NoError(t, err)

	assert.True(t, result.DryRun)
	assert.Equal(t, 2, len(result.Operations))
	assert.Equal(t, 2, result.Summary.Planned)
	assert.Equal(t, 0, result.Summary.Renamed)

	// Ensure determinism by checking order
	assert.Equal(t, f1, result.Operations[0].Source)
	assert.Equal(t, filepath.Join(tempDir, "file_one_renamed.txt"), result.Operations[0].Destination)

	assert.Equal(t, f2, result.Operations[1].Source)
	assert.Equal(t, filepath.Join(tempDir, "file_two_renamed.txt"), result.Operations[1].Destination)
}
