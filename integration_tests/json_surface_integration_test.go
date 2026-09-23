package integration_tests

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arran4/rntocase"
)

func runCLI(t *testing.T, binPath string, args ...string) (string, string, error) {
	cmd := exec.Command(binPath, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	require.NoError(t, err)

	stderrPipe, err := cmd.StderrPipe()
	require.NoError(t, err)

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	stdoutBytes, err := io.ReadAll(stdoutPipe)
	require.NoError(t, err)

	stderrBytes, err := io.ReadAll(stderrPipe)
	require.NoError(t, err)

	err = cmd.Wait()
	return string(stdoutBytes), string(stderrBytes), err
}

func TestCommandJSONSurface(t *testing.T) {
	binPath := sharedBinPath

	t.Run("basic json output structure", func(t *testing.T) {
		tempDir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_one.txt"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_two.txt"), []byte("test"), 0644))

		stdout, stderr, err := runCLI(t, binPath, "camel", "--json", filepath.Join(tempDir, "file_one.txt"), filepath.Join(tempDir, "file_two.txt"))
		require.NoError(t, err)
		assert.Empty(t, stderr)

		var output struct {
			Operations []rntocase.RenameOperation `json:"operations"`
			Summary    struct {
				Renamed int `json:"renamed"`
			} `json:"summary"`
		}

		err = json.Unmarshal([]byte(stdout), &output)
		require.NoError(t, err)

		assert.Len(t, output.Operations, 2)
		assert.Equal(t, 2, output.Summary.Renamed)
	})

	t.Run("json with preflight error", func(t *testing.T) {
		tempDir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "fileOne.txt"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_one.txt"), []byte("test"), 0644))

		stdout, stderr, err := runCLI(t, binPath, "camel", "--json", filepath.Join(tempDir, "fileOne.txt"), filepath.Join(tempDir, "file_one.txt"))

		require.Error(t, err) // Exit code should be non-zero
		assert.Contains(t, stderr, "collision")

		var output struct {
			Operations []rntocase.RenameOperation `json:"operations"`
			Summary    struct {
				Renamed int `json:"renamed"`
			} `json:"summary"`
		}

		err = json.Unmarshal([]byte(stdout), &output)
		require.NoError(t, err)
		assert.Len(t, output.Operations, 2)
		assert.Equal(t, 0, output.Summary.Renamed)
	})

	t.Run("json with mutually exclusive interactive mode", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test"), 0644))

		_, stderr, err := runCLI(t, binPath, "camel", "--json", "--interactive", filepath.Join(tempDir, "file1.txt"))
		require.Error(t, err)
		assert.Contains(t, stderr, "cannot use interactive mode with JSON output")
	})

	t.Run("json ignores log output on stdout", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "file_1.txt"), []byte("test"), 0644))

		stdout, stderr, err := runCLI(t, binPath, "camel", "--json", "--dry-run", filepath.Join(tempDir, "file_1.txt"))
		require.NoError(t, err)

		// JSON should be parsable
		var output interface{}
		err = json.Unmarshal([]byte(stdout), &output)
		require.NoError(t, err)

		// Verify there is no JSON document output to stderr
		assert.NotContains(t, stderr, "{")
	})

	t.Run("recursive explicit include/exclude with json", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "src", "pkg"), 0755))

		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "src", "file_1.go"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "src", "pkg", "file_2.go"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "src", "pkg", "ignore_me.txt"), []byte("test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(tempDir, "src", "pkg", "vendor.go"), []byte("test"), 0644))

		stdout, stderr, err := runCLI(t, binPath, "camel", "--json", "-R", "--include", "**/*.go", "--exclude", "**/vendor.go", tempDir)
		require.NoError(t, err)
		assert.Empty(t, stderr)

		var output struct {
			Operations []rntocase.RenameOperation `json:"operations"`
			Summary    struct {
				Renamed int `json:"renamed"`
			} `json:"summary"`
		}

		err = json.Unmarshal([]byte(stdout), &output)
		require.NoError(t, err)

		assert.Equal(t, 2, output.Summary.Renamed)

		foundFile1 := false
		foundFile2 := false
		for _, op := range output.Operations {
			if strings.Contains(op.Source, "file_1.go") {
				foundFile1 = true
			}
			if strings.Contains(op.Source, "file_2.go") {
				foundFile2 = true
			}
		}
		assert.True(t, foundFile1)
		assert.True(t, foundFile2)
	})
}
