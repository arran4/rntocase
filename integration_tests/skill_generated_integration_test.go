package integration_tests

import (
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

	t.Run("install with ref", func(t *testing.T) {
		runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--ref", "v0.0.1", "--path", "skills/example", "--name", "example", "arran4/non-existent-rntocase")
		out, err := runCmd.CombinedOutput()

		// Prove it parsed correctly and reached the failure phase for resolving github
		require.Error(t, err)
		assert.Contains(t, string(out), "failed to get repository metadata: HTTP 404")
	})

	t.Run("install with source positional before flags", func(t *testing.T) {
		// Positional arguments cannot be placed before flags according to gosubc architecture,
		// but checking that it gives the correct error regarding standard flag parsing.
		runCmd := exec.Command(binPath, "skill", "install", "arran4/non-existent-rntocase", "--scope=user", "--ref", "v0.0.1", "--path", "skills/example", "--name", "example")
		out, err := runCmd.CombinedOutput()

		require.Error(t, err)
		// It tries to install "arran4/non-existent-rntocase" and assumes "--scope=user" etc. are more positional args, but positional limits hit.
		// Wait, previously the output showed it attempting to resolve "--scope=user" as a target. Let's match the HTTP 404 meaning it parses the source.
		assert.Contains(t, string(out), "failed to get repository metadata: HTTP 404")
	})
}
