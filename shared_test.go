package rntocase

import (
	"fmt"
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

func TestRunSafeBatch(t *testing.T) {
	// Success case
	values := []string{"a", "b"}
	algoSuccess := func(s string) (string, error) { return s + "!", nil }
	var results []string
	runSafeBatch(algoSuccess, values, func(s string) bool {
		results = append(results, s)
		return true
	})
	if len(results) != 2 || results[0] != "a!" || results[1] != "b!" {
		t.Errorf("Unexpected results for success case: %v", results)
	}

	// Error case
	results = nil
	algoError := func(s string) (string, error) {
		if s == "a" {
			return "", fmt.Errorf("error")
		}
		return s + "!", nil
	}
	runSafeBatch(algoError, values, func(s string) bool {
		results = append(results, s)
		return true
	})
	if len(results) != 2 || results[0] != "!!!Error!!!" || results[1] != "b!" {
		t.Errorf("Unexpected results for error case: %v", results)
	}

	// Panic case
	results = nil
	algoPanic := func(s string) (string, error) {
		if s == "a" {
			panic("oops")
		}
		return s + "!", nil
	}
	runSafeBatch(algoPanic, values, func(s string) bool {
		results = append(results, s)
		return true
	})
	if len(results) != 2 || results[0] != "!!!Error!!!" || results[1] != "b!" {
		t.Errorf("Unexpected results for panic case: %v", results)
	}

	// Yield interrupt case
	results = nil
	runSafeBatch(algoSuccess, values, func(s string) bool {
		results = append(results, s)
		return false // Stop
	})
	if len(results) != 1 || results[0] != "a!" {
		t.Errorf("Unexpected results for interrupt case: %v", results)
	}
}

func TestRenderUsageTable(t *testing.T) {
	algos := map[string]func(string) (string, error){
		"echo": func(s string) (string, error) { return s, nil },
		"panic": func(s string) (string, error) {
			if strings.Contains(s, "Hello") {
				panic("oops")
			}
			return s, nil
		},
	}
	// It writes to stderr, so we might see output in test logs, but it shouldn't crash.
	// We can't easily capture output here without redefining ExampleGroups or mocking things.
	// But running it ensures no panic propagates out.
	RenderUsageTable(algos)
}
