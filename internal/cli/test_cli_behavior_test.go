package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arran4/rntocase/internal/skill"
	"github.com/stretchr/testify/assert"
)

// Because skill target paths are complex, we override HOME for the test duration
func setupMockHome(t *testing.T) string {
	t.Helper()
	homeDir, err := os.MkdirTemp("", "mock-home-*")
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(homeDir)
	})
	t.Setenv("HOME", homeDir)
	return homeDir
}

func TestRunSkillInstall_ReplaceFlagFailureLeavesPriorIntact(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed skill under user scope for 'common' agent
	// Path should be $HOME/.agents/skills/my-bad-skill
	destDir := filepath.Join(homeDir, ".agents", "skills", "my-bad-skill")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("original working version"), 0644)
	assert.NoError(t, err)

	// Save valid metadata so it's "managed"
	err = skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "my-bad-skill",
		OriginalSource: "local/mock",
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
	})
	assert.NoError(t, err)

	// Try installing a local dir over it WITH --replace, but the local dir lacks SKILL.md
	// so the validation phase inside ReplaceSafely will fail.
	sourceDir, err := os.MkdirTemp("", "bad-source-*")
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(sourceDir)
	})
	// We do NOT write SKILL.md to sourceDir

	// Run command
	err = RunSkillInstall([]string{"--replace", "--scope=user", "--agent=common", sourceDir, "my-bad-skill"})

	// Expect failure because of missing SKILL.md
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "skill must contain a SKILL.md file")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "original working version", string(content))
}

func TestUpdateSingleSkill_EmbeddedValidationFailureLeavesPriorIntact(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed "official" skill
	destDir := filepath.Join(homeDir, ".agents", "skills", "official-skill")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("original working version"), 0644)
	assert.NoError(t, err)

	meta := &skill.Metadata{
		Name:           "official-skill",
		OriginalSource: "official", // Triggers the embedded path
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
		ContentDigest:  "some-old-digest", // Provide a digest so update isn't skipped for being identical
	}
	err = skill.SaveMetadata(destDir, meta)
	assert.NoError(t, err)

	// Inject failure in ExtractEmbeddedSkill via the explicit dependency injection parameter
	mockExtract := func(skillName, destDir string) error {
		// Return success but do not create SKILL.md to simulate validation failure
		return nil
	}

	meta.OriginalSource = "official" // triggers embedded logic
	err = updateSingleSkillWithExtractor("official-skill", meta, destDir, true, mockExtract)

	// Expect it to fail
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "new skill version must contain a SKILL.md file")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "original working version", string(content))
}

func TestUpdateSingleSkill_EmbeddedExtractionFailureLeavesPriorIntact(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed "official" skill
	destDir := filepath.Join(homeDir, ".agents", "skills", "official-skill")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("original working version"), 0644)
	assert.NoError(t, err)

	meta := &skill.Metadata{
		Name:           "official-skill",
		OriginalSource: "official",
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
		ContentDigest:  "some-old-digest",
	}
	err = skill.SaveMetadata(destDir, meta)
	assert.NoError(t, err)

	// Inject extraction failure via explicit dependency injection
	mockExtract := func(skillName, destDir string) error {
		return fmt.Errorf("simulated embedded extraction failure")
	}

	err = updateSingleSkillWithExtractor("official-skill", meta, destDir, true, mockExtract)

	// Expect it to fail
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "simulated embedded extraction failure")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "original working version", string(content))
}
