package rntocase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenameFiles(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "rntocase_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create dummy files
	files := []string{
		"Hello World.txt",
		"testFile.go",
		"AlreadyLower.txt",
	}

	for _, f := range files {
		path := filepath.Join(tempDir, f)
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Rename function (to lowercase)
	renameFunc := func(s string) (string, error) {
		return strings.ToLower(s), nil
	}

	// Paths to rename
	var paths []string
	for _, f := range files {
		paths = append(paths, filepath.Join(tempDir, f))
	}

	// Run RenameFiles
	if err := RenameFiles(paths, renameFunc, false, false); err != nil {
		t.Fatalf("RenameFiles failed: %v", err)
	}

	// Verify results
	expectedFiles := []string{
		"hello world.txt",
		"testfile.go",
		"alreadylower.txt",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(tempDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", path)
		}
	}

	// Verify original files are gone
	// "AlreadyLower.txt" -> "alreadylower.txt". strings.ToLower("AlreadyLower") is "alreadylower".
	// So it changes.

	// Check if original "Hello World.txt" is gone
	if _, err := os.Stat(filepath.Join(tempDir, "Hello World.txt")); err == nil {
		t.Errorf("Original file 'Hello World.txt' should not exist")
	}
}
