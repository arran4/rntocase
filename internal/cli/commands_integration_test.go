package cli

import (
	"io"
	"strings"

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

	err := RunCamel(false, false, false, filePath)
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

	err = RunPascal(false, false, false, filePath2)
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
	err := RunDot("", false, false, false, filePath)
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
	err := RunTrim("_", false, false, false, filePath)
	if err != nil {
		t.Fatalf("RunTrim failed: %v", err)
	}

	expectedTrimPath := filepath.Join(tmpDir, "hello_world.txt")
	if _, err := os.Stat(expectedTrimPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", expectedTrimPath)
	}
}

func TestCommandJSONSurface(t *testing.T) {
	tempDir := t.TempDir()

	// Helper to capture stdout/stderr separately
	captureOutput := func(f func() error) (string, string, error) {
		oldStdout := os.Stdout
		oldStderr := os.Stderr

		rOut, wOut, _ := os.Pipe()
		rErr, wErr, _ := os.Pipe()

		os.Stdout = wOut
		os.Stderr = wErr

		errCh := make(chan error, 3)

		go func() {
			errCh <- f()
			errCh <- wOut.Close()
			errCh <- wErr.Close()
		}()

		var bufOut strings.Builder
		var bufErr strings.Builder

		_, _ = io.Copy(&bufOut, rOut)
		_, _ = io.Copy(&bufErr, rErr)

		os.Stdout = oldStdout
		os.Stderr = oldStderr

		_ = rOut.Close()
		_ = rErr.Close()

		execErr := <-errCh
		<-errCh
		<-errCh

		return bufOut.String(), bufErr.String(), execErr
	}

	t.Run("interactive rejected", func(t *testing.T) {
		_, _, err := captureOutput(func() error {
			return RunCamel(true, false, true, "test")
		})
		if err == nil || !strings.Contains(err.Error(), "cannot use interactive mode with JSON output") {
			t.Errorf("Expected interactive rejection error, got %v", err)
		}
	})

	t.Run("dry run JSON", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "dry_run_camel.txt")
		if err := os.WriteFile(fPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		stdout, stderr, err := captureOutput(func() error {
			return RunCamel(true, true, false, fPath)
		})

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if stderr != "" {
			t.Errorf("Expected empty stderr, got %s", stderr)
		}

		// Assert valid JSON
		if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
			t.Errorf("Expected stdout to be JSON, got %s", stdout)
		}
	})

	t.Run("execution JSON and failure", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "exec_fail_camel.txt")
		// file does not exist, rename should fail (or fail planning)

		stdout, stderr, err := captureOutput(func() error {
			return RunCamel(true, false, false, fPath)
		})

		if err == nil {
			t.Error("Expected error due to missing file")
		}

		// Error should still result in parseable JSON in stdout
		if !strings.HasPrefix(strings.TrimSpace(stdout), "{") {
			t.Errorf("Expected stdout to be JSON even on error, got %s", stdout)
		}

		// Diagnostics might be in stderr or just returned as error
		_ = stderr // stderr can be anything
	})
}
