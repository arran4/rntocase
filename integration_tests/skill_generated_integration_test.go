package integration_tests

import (
	"archive/tar"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedCommand_SkillInstall_Integration(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "rntocase")
	cmd := exec.Command("go", "build", "-o", binPath, "github.com/arran4/rntocase/cmd/rntocase")
	cmd.Dir = "../"
	err := cmd.Run()
	require.NoError(t, err, "failed to build CLI")

	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: generated-skill\ndescription: desc\n---"), 0644))

	t.Run("basic user scope", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--name=generated-skill", sourceDir)
		out, err := runCmd.CombinedOutput()
		require.NoError(t, err, "failed to run command: %s", string(out))

		destDir := filepath.Join(homeDir, ".agents", "skills", "generated-skill")
		content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
		require.NoError(t, err, "could not read installed SKILL.md")
		assert.Equal(t, "---\nname: generated-skill\ndescription: desc\n---", string(content))
	})

	t.Run("bundled official skill", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "rntocase")
		out, err := runCmd.CombinedOutput()
		require.NoError(t, err, "failed to run command: %s", string(out))

		destDir := filepath.Join(homeDir, ".agents", "skills", "rntocase")
		content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
		require.NoError(t, err, "could not read installed SKILL.md")
		assert.Contains(t, string(content), "name: rntocase")
	})

	t.Run("update explicit form", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "update", "--scope=user", "rntocase")
		out, err := runCmd.CombinedOutput()
		require.NoError(t, err, "failed to run command: %s", string(out))
		assert.Contains(t, string(out), "updated")
	})

	t.Run("update non existing form", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "update", "--scope=user", "non-existent")
		err := runCmd.Run()
		require.Error(t, err, "expected error")
	})

	var requestedCommitsPath, requestedTarballPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/arran4/mock-rntocase/commits/v0.0.1" {
			requestedCommitsPath = r.URL.Path
			w.WriteHeader(http.StatusOK)
			// Return a minimal JSON response mimicking a commit payload
			_, _ = w.Write([]byte(`{"sha": "mock-sha-12345"}`))
			return
		}
		if r.URL.Path == "/repos/arran4/mock-rntocase/tarball/mock-sha-12345" {
			requestedTarballPath = r.URL.Path
			w.WriteHeader(http.StatusOK)
			gw := gzip.NewWriter(w)
			tw := tar.NewWriter(gw)

			// Directory entry
			_ = tw.WriteHeader(&tar.Header{
				Name:     "arran4-mock-rntocase-12345/",
				Typeflag: tar.TypeDir,
				Mode:     0755,
			})

			// Subdirectory entry
			_ = tw.WriteHeader(&tar.Header{
				Name:     "arran4-mock-rntocase-12345/skills/example/",
				Typeflag: tar.TypeDir,
				Mode:     0755,
			})

			// SKILL.md file
			content := "---\nname: example\ndescription: mock skill desc\n---\n"
			_ = tw.WriteHeader(&tar.Header{
				Name: "arran4-mock-rntocase-12345/skills/example/SKILL.md",
				Mode: 0644,
				Size: int64(len(content)),
			})
			_, _ = tw.Write([]byte(content))

			_ = tw.Close()
			_ = gw.Close()
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	t.Run("install with ref using mock network", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--ref", "v0.0.1", "--path", "skills/example", "--name", "example", "arran4/mock-rntocase")
		runCmd.Env = append(os.Environ(), "RNTOCASE_GITHUB_API_URL="+ts.URL)

		out, err := runCmd.CombinedOutput()
		require.NoError(t, err, "failed to run command: %s", string(out))

		destDir := filepath.Join(homeDir, ".agents", "skills", "example")

		// Check that SKILL.md was extracted
		content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
		require.NoError(t, err, "could not read installed SKILL.md")
		assert.Equal(t, "---\nname: example\ndescription: mock skill desc\n---\n", string(content))

		// Check the metadata sidecar
		metaContent, err := os.ReadFile(filepath.Join(destDir, ".rntocase-skill.json"))
		require.NoError(t, err, "could not read metadata sidecar")

		assert.Contains(t, string(metaContent), `"name": "example"`)
		assert.Contains(t, string(metaContent), `"owner_repo": "arran4/mock-rntocase"`)
		assert.Contains(t, string(metaContent), `"requested_ref": "v0.0.1"`)
		assert.Contains(t, string(metaContent), `"source_revision": "mock-sha-12345"`)
		assert.Contains(t, string(metaContent), `"path_within": "skills/example"`)

		// Assert strict route hits
		assert.Equal(t, "/repos/arran4/mock-rntocase/commits/v0.0.1", requestedCommitsPath)
		assert.Equal(t, "/repos/arran4/mock-rntocase/tarball/mock-sha-12345", requestedTarballPath)
	})

	t.Run("install with ref positional after flags (documented form)", func(t *testing.T) {
		// Documented form includes flags explicitly before positional:
		// rntocase skill install --ref v1.2.0 --path skills/example --name example owner/repo
		// Wait, the prior run tested `skill install --scope=user --ref v0.0.1 ... arran4/mock-rntocase` which IS flags before positional.
		// The manager asked to keep "the documented supported flag/positional ordering explicitly exercised".
		// We will test `install --ref v0.0.1 --path skills/example --name example arran4/mock-rntocase` to explicitly match documentation.

		// The previous mock `requestedCommitsPath` state is preserved in the singleton test server, we reset it or just make a new call
		requestedCommitsPath = ""
		requestedTarballPath = ""

		// Passing --replace to overwrite the existing "example" directory installed in the previous step
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--replace", "--ref", "v0.0.1", "--path", "skills/example", "--name", "example", "arran4/mock-rntocase")
		runCmd.Env = append(os.Environ(), "RNTOCASE_GITHUB_API_URL="+ts.URL)

		out, err := runCmd.CombinedOutput()
		require.NoError(t, err, "failed to run command with flags before positional (documented form): %s", string(out))

		destDir := filepath.Join(homeDir, ".agents", "skills", "example")
		content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
		require.NoError(t, err, "could not read installed SKILL.md")
		assert.Equal(t, "---\nname: example\ndescription: mock skill desc\n---\n", string(content))

		assert.Equal(t, "/repos/arran4/mock-rntocase/commits/v0.0.1", requestedCommitsPath)
		assert.Equal(t, "/repos/arran4/mock-rntocase/tarball/mock-sha-12345", requestedTarballPath)
	})

	t.Run("install with ref returns 404 from mock network", func(t *testing.T) {
		// Use a specific path check that will explicitly hit the 404 case.
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--ref", "v0.0.2", "--path", "skills/example", "--name", "example3", "arran4/mock-rntocase-notfound")
		runCmd.Env = append(os.Environ(), "RNTOCASE_GITHUB_API_URL="+ts.URL)

		out, err := runCmd.CombinedOutput()
		require.Error(t, err, "expected error when github returns 404")
		assert.Contains(t, string(out), "failed to get repository metadata: HTTP 404")
	})

}
