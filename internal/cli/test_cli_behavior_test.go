package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arran4/rntocase/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Because skill target paths are complex, we override HOME for the test duration
func setupMockHome(t *testing.T) string {
	t.Helper()
	homeDir := t.TempDir()
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

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: original-skill\ndescription: test skill\n---"), 0644)
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
	sourceDir := t.TempDir()
	// We do NOT write SKILL.md to sourceDir

	// Run command
	err = RunSkillInstall("user", "common", true, "", "", "", sourceDir, "my-bad-skill")

	// Expect failure because of missing SKILL.md
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "skill must contain a SKILL.md file")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "---\nname: original-skill\ndescription: test skill\n---", string(content))
}

func TestRunSkillInstall_RejectsMalformedFrontmatterBeforeReplacingManagedSkill(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "my-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	previousSkill := []byte("---\nname: my-skill\ndescription: previous\n---\n")
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), previousSkill, 0644))
	previousMetadata := &skill.Metadata{Name: "my-skill", OriginalSource: "local/mock", InstallerApp: "rntocase"}
	require.NoError(t, skill.SaveMetadata(destDir, previousMetadata))
	previousMetadataBytes, err := os.ReadFile(filepath.Join(destDir, skill.MetadataFileName))
	require.NoError(t, err)

	sourceDir := t.TempDir()
	malformedSkill := []byte("---\nname: my-skill\ndescription: invalid\n---not-a-delimiter\n")
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), malformedSkill, 0644))

	err = RunSkillInstall("user", "common", true, "", "", "my-skill", sourceDir, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing closing '---' on a standalone line")

	actualSkill, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, previousSkill, actualSkill)
	actualMetadata, err := os.ReadFile(filepath.Join(destDir, skill.MetadataFileName))
	require.NoError(t, err)
	assert.Equal(t, previousMetadataBytes, actualMetadata)
}

func TestUpdateSingleSkill_RejectsMalformedFrontmatterBeforeReplacingManagedSkill(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "official-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	previousSkill := []byte("---\nname: official-skill\ndescription: previous\n---\n")
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), previousSkill, 0644))
	meta := &skill.Metadata{Name: "official-skill", OriginalSource: "official", InstallerApp: "rntocase", ContentDigest: "old-digest"}
	require.NoError(t, skill.SaveMetadata(destDir, meta))
	previousMetadata, err := os.ReadFile(filepath.Join(destDir, skill.MetadataFileName))
	require.NoError(t, err)

	malformedExtractor := func(_ string, stagingDir string) error {
		return os.WriteFile(filepath.Join(stagingDir, "SKILL.md"), []byte("---\nname: official-skill\ndescription: invalid\n---not-a-delimiter\n"), 0644)
	}
	_, err = updateSingleSkillWithExtractor("official-skill", meta, destDir, true, malformedExtractor)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing closing '---' on a standalone line")

	actualSkill, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, previousSkill, actualSkill)
	actualMetadata, err := os.ReadFile(filepath.Join(destDir, skill.MetadataFileName))
	require.NoError(t, err)
	assert.Equal(t, previousMetadata, actualMetadata)
}

func TestUpdateSingleSkill_EmbeddedValidationFailureLeavesPriorIntact(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed "official" skill
	destDir := filepath.Join(homeDir, ".agents", "skills", "official-skill")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: original-skill\ndescription: test skill\n---"), 0644)
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
	mockExtract := func(assetName, destDir string) error {
		// Return success but do not create SKILL.md to simulate validation failure
		return nil
	}

	meta.OriginalSource = "official" // triggers embedded logic
	_, err = updateSingleSkillWithExtractor("official-skill", meta, destDir, true, mockExtract)

	// Expect it to fail
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "new skill version must contain a SKILL.md file")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "---\nname: original-skill\ndescription: test skill\n---", string(content))
}

func TestUpdateSingleSkill_EmbeddedExtractionFailureLeavesPriorIntact(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed "official" skill
	destDir := filepath.Join(homeDir, ".agents", "skills", "official-skill")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: original-skill\ndescription: test skill\n---"), 0644)
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
	mockExtract := func(assetName, destDir string) error {
		return fmt.Errorf("simulated embedded extraction failure")
	}

	_, err = updateSingleSkillWithExtractor("official-skill", meta, destDir, true, mockExtract)

	// Expect it to fail
	if err == nil {
		t.Fatalf("Expected an error but got nil. Err: %v", err)
	}
	assert.Contains(t, err.Error(), "simulated embedded extraction failure")

	// Ensure prior working installation is entirely intact
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "---\nname: original-skill\ndescription: test skill\n---", string(content))
}

// Regression test for #38: official skill installed under a custom name still targets "rntocase" asset
func TestUpdateSingleSkill_CustomDestNamePreservesOfficialAsset(t *testing.T) {
	homeDir := setupMockHome(t)

	// Pre-create an installed, managed "official" skill under a CUSTOM name
	destDir := filepath.Join(homeDir, ".agents", "skills", "my-custom-skill-name")
	err := os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: original-skill\ndescription: test skill\n---"), 0644)
	assert.NoError(t, err)

	meta := &skill.Metadata{
		Name:           "my-custom-skill-name", // Custom name
		OriginalSource: "official",
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
		ContentDigest:  "some-old-digest",
	}
	err = skill.SaveMetadata(destDir, meta)
	assert.NoError(t, err)

	capturedExtractedName := ""
	mockExtract := func(assetName, destDir string) error {
		capturedExtractedName = assetName

		// Create SKILL.md to pass validation
		return os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: my-custom-skill-name\ndescription: test skill\n---"), 0644)
	}

	_, err = updateSingleSkillWithExtractor("my-custom-skill-name", meta, destDir, true, mockExtract)
	assert.NoError(t, err)

	// Confirm the production integration invokes the right asset name by design.
	assert.Equal(t, "rntocase", capturedExtractedName)
}
