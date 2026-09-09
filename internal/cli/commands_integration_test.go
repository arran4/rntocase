package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntegrationCamelVsPascal(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "hello world.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	err := RunCamel(false, false, filePath)
	if err != nil {
		t.Fatalf("RunCamel failed: %v", err)
	}

	expectedCamelPath := filepath.Join(tmpDir, "helloWorld.txt")
	if _, err := os.Stat(expectedCamelPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedCamelPath)
	}

	// Test Pascal
	filePath2 := filepath.Join(tmpDir, "hello world 2.txt")
	if err := os.WriteFile(filePath2, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	err = RunPascal(false, false, filePath2)
	if err != nil {
		t.Fatalf("RunPascal failed: %v", err)
	}

	expectedPascalPath := filepath.Join(tmpDir, "HelloWorld2.txt")
	if _, err := os.Stat(expectedPascalPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedPascalPath)
	}
}

func TestIntegrationDotDefault(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "hello world.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Call RunDot with empty delimiter to test default
	err := RunDot("", false, false, filePath)
	if err != nil {
		t.Fatalf("RunDot failed: %v", err)
	}

	expectedDotPath := filepath.Join(tmpDir, "hello.world.txt")
	if _, err := os.Stat(expectedDotPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedDotPath)
	}
}

func TestIntegrationTrim(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "_hello_world_.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	// Call RunTrim
	err := RunTrim("_", false, false, filePath)
	if err != nil {
		t.Fatalf("RunTrim failed: %v", err)
	}

	expectedTrimPath := filepath.Join(tmpDir, "hello_world.txt")
	if _, err := os.Stat(expectedTrimPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedTrimPath)
	}
}
