package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

// Tests that used to be integration tests but don't need real subprocess isolation.
func TestSkillInstall_BasicLocal(t *testing.T) {
	homeDir := setupMockHome(t)
	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: generated-skill\ndescription: desc\n---"), 0644))

	err := RunSkillInstall("user", "", false, "", "", "generated-skill", sourceDir, "")
	require.NoError(t, err)

	destDir := filepath.Join(homeDir, ".agents", "skills", "generated-skill")
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, "---\nname: generated-skill\ndescription: desc\n---", string(content))
}

func TestSkillInstall_BundledOfficial(t *testing.T) {
	homeDir := setupMockHome(t)

	err := RunSkillInstall("user", "", false, "", "", "", "rntocase", "")
	require.NoError(t, err)

	destDir := filepath.Join(homeDir, ".agents", "skills", "rntocase")
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: rntocase")
}

func TestSkillUpdate_ExplicitForm(t *testing.T) {
	_ = setupMockHome(t)
	err := RunSkillInstall("user", "", false, "", "", "", "rntocase", "")
	require.NoError(t, err)

	err = RunSkillUpdate("user", "", false, false, "rntocase")
	require.NoError(t, err)
}

func TestSkillUpdate_NonExistingForm(t *testing.T) {
	_ = setupMockHome(t)
	err := RunSkillUpdate("user", "", false, false, "non-existent")
	require.Error(t, err)
}
