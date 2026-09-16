package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: generated-skill\ndescription: desc\n---"), 0644))

	runCmd := exec.Command(binPath, "skill", "install", "--scope=user", "--name=generated-skill", sourceDir)
	out, err := runCmd.CombinedOutput()
	require.NoError(t, err, "failed to run command: %s", string(out))

	destDir := filepath.Join(homeDir, ".agents", "skills", "generated-skill")
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err, "could not read installed SKILL.md")
	assert.Equal(t, "---\nname: generated-skill\ndescription: desc\n---", string(content))
}
