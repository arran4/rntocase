package main

import (
	"encoding/json"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/arran4/rntocase/cmd"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, f func() error) (string, string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, errOut := os.Pipe()
	if errOut != nil {
		t.Fatalf("Failed to create stdout pipe: %v", errOut)
	}
	rErr, wErr, errErr := os.Pipe()
	if errErr != nil {
		t.Fatalf("Failed to create stderr pipe: %v", errErr)
	}

	os.Stdout = wOut
	os.Stderr = wErr

	errCh := make(chan error, 1)
	outCh := make(chan string)
	errStrCh := make(chan string)

	go func() {
		errCh <- f()
		if err := wOut.Close(); err != nil {
			t.Errorf("Error closing wOut: %v", err)
		}
		if err := wErr.Close(); err != nil {
			t.Errorf("Error closing wErr: %v", err)
		}
	}()

	go func() {
		var bufOut strings.Builder
		if _, err := io.Copy(&bufOut, rOut); err != nil {
			t.Errorf("Error copying stdout: %v", err)
		}
		outCh <- bufOut.String()
	}()

	go func() {
		var bufErr strings.Builder
		if _, err := io.Copy(&bufErr, rErr); err != nil {
			t.Errorf("Error copying stderr: %v", err)
		}
		errStrCh <- bufErr.String()
	}()

	stdoutStr := <-outCh
	stderrStr := <-errStrCh
	execErr := <-errCh

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if err := rOut.Close(); err != nil {
		t.Errorf("Error closing rOut: %v", err)
	}
	if err := rErr.Close(); err != nil {
		t.Errorf("Error closing rErr: %v", err)
	}

	return stdoutStr, stderrStr, execErr
}

func TestCommandJSONSurface(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("interactive rejected", func(t *testing.T) {
		_, _, err := captureOutput(t, func() error {
			root, err := NewRoot("rntocase", "dev", "none", "unknown")
			if err != nil {
				return err
			}
			err = root.Execute([]string{"camel", "--json", "--interactive", "test"})
			if err != nil {
				if e, ok := err.(*cmd.ErrExitCode); ok {
					if e.Err != nil {
						return e.Err
					}
				}
				return err
			}
			return nil
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

		wrapper := func() error {
			root, err := NewRoot("rntocase", "dev", "none", "unknown")
			if err != nil {
				return err
			}
			err = root.Execute([]string{"camel", "--json", "--dry-run", fPath})
			if err != nil {
				if e, ok := err.(*cmd.ErrExitCode); ok {
					return e.Err
				}
				return err
			}
			return nil
		}

		stdout, stderr, err := captureOutput(t, wrapper)

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if stderr != "" {
			t.Errorf("Expected empty stderr, got %s", stderr)
		}

		// Assert strict JSON via Unmarshal
		var res rntocase.RenameResult
		if err := json.Unmarshal([]byte(stdout), &res); err != nil {
			t.Fatalf("Failed to decode stdout JSON: %v, output: %s", err, stdout)
		}
		if !res.DryRun {
			t.Errorf("Expected DryRun true in JSON result")
		}
		if res.Summary.Planned != 1 {
			t.Errorf("Expected 1 planned, got %d", res.Summary.Planned)
		}
	})

	t.Run("execution JSON and failure", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "exec_fail_camel.txt")
		// file does not exist, rename should fail (or fail planning)

		wrapper := func() error {
			root, err := NewRoot("rntocase", "dev", "none", "unknown")
			if err != nil {
				return err
			}
			err = root.Execute([]string{"camel", "--json", fPath})
			if err != nil {
				if e, ok := err.(*cmd.ErrExitCode); ok {
					if e.Err != nil {
						return e.Err
					}
					if e.Code != 0 {
						return fmt.Errorf("exit code %d", e.Code)
					}
				}
				return err
			}
			return nil
		}

		stdout, stderr, err := captureOutput(t, wrapper)

		if err == nil {
			t.Error("Expected error due to missing file")
		}

		// Error should still result in parseable JSON in stdout
		var res rntocase.RenameResult
		if err := json.Unmarshal([]byte(stdout), &res); err != nil {
			t.Fatalf("Failed to decode stdout JSON on error: %v, output: %s", err, stdout)
		}
		if res.Summary.Failed != 1 {
			t.Errorf("Expected 1 failed operation, got %d", res.Summary.Failed)
		}

		// Diagnostics might be in stderr or just returned as error
		_ = stderr // stderr can be anything
	})

	t.Run("execution JSON success", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "exec_success_camel.txt")
		if err := os.WriteFile(fPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		wrapper := func() error {
			root, err := NewRoot("rntocase", "dev", "none", "unknown")
			if err != nil {
				return err
			}
			err = root.Execute([]string{"camel", "--json", fPath})
			if err != nil {
				if e, ok := err.(*cmd.ErrExitCode); ok {
					return e.Err
				}
				return err
			}
			return nil
		}

		stdout, _, err := captureOutput(t, wrapper)

		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		var res rntocase.RenameResult
		if err := json.Unmarshal([]byte(stdout), &res); err != nil {
			t.Fatalf("Failed to decode stdout JSON: %v, output: %s", err, stdout)
		}
		if res.Summary.Renamed != 1 {
			t.Errorf("Expected 1 renamed, got %d", res.Summary.Renamed)
		}
	})
}
