package main

import (
	"os/exec"
	"path/filepath"
	"os"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestGeneratedCommand_SkillInstall_Integration(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "rntocase")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "failed to build CLI")

	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("ok"), 0644))

	runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--name=generated-skill", sourceDir)
	out, err := runCmd.CombinedOutput()
	require.NoError(t, err, "failed to run command: %s", string(out))

	destDir := filepath.Join(homeDir, ".agents", "skills", "generated-skill")
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err, "could not read installed SKILL.md")
	assert.Equal(t, "ok", string(content))
}
