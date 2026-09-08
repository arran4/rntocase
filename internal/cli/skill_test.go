package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/arran4/rntocase/internal/skill"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSkill_RequiresSubcommand(t *testing.T) {
	err := RunSkill([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a subcommand")
}

func TestRunSkill_UnknownSubcommand(t *testing.T) {
	err := RunSkill([]string{"unknown_cmd"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown skill subcommand: unknown_cmd")
}

func TestRunSkillInstall_RequiresSource(t *testing.T) {
	err := RunSkillInstall([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill install [--replace] <source>")
}

func TestRunSkillUpdate_RequiresName(t *testing.T) {
	err := RunSkillUpdate([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill update <name> or skill update --all")
}

func TestRunSkillRemove_RequiresName(t *testing.T) {
	err := RunSkillRemove([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill remove <name>")
}

func TestRunSkillInspect_RequiresName(t *testing.T) {
	err := RunSkillInspect([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill inspect <name>")
}

func TestRunSkillUpdate_LocalSkillError(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("ok"), 0644))

	require.NoError(t, skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "local-skill1",
		OriginalSource: "local",
	}))

	err := RunSkillUpdate([]string{"--scope", "user", "local-skill1"})
	assert.Error(t, err, "Updating a directly targeted local skill should preserve previous single-skill error semantics")
	assert.Contains(t, err.Error(), "is locally installed and cannot be updated automatically")
}

func TestRunSkillUpdate_InspectionFailure(t *testing.T) {
	err := RunSkillUpdate([]string{"--scope", "user", "non-existent-skill"})
	assert.Error(t, err, "Should fail inspection")
	assert.Contains(t, err.Error(), "skill 'non-existent-skill' not found")
}

func TestRunSkillUpdate_All_InspectionFailure(t *testing.T) {
	homeDir := setupMockHome(t)

	// Valid skill
	destDir1 := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "local-skill1", OriginalSource: "local"}))

	// Second skill that ListInstalledSkills can find, but InspectSkill will fail on.
	// This happens if the directory name doesn't match the inner skill name.
	destDir2 := filepath.Join(homeDir, ".agents", "skills", "mismatched-dir")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "inner-name", OriginalSource: "local"}))

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	// Mismatched dir will cause InspectSkill to look for "inner-name" directory, which doesn't exist under .agents/skills
	assert.Contains(t, err.Error(), "inner-name")
	assert.Contains(t, err.Error(), "inspection failed")
	assert.NotContains(t, err.Error(), "local-skill1")
}

func TestRunSkillUpdate_PartialFailure(t *testing.T) {
	homeDir := setupMockHome(t)

	destDir1 := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "local-skill1", OriginalSource: "local"}))

	destDir2 := filepath.Join(homeDir, ".agents", "skills", "broken-skill")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "broken-skill", OwnerRepo: "dummy/repo"}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "broken-skill")
	assert.Contains(t, err.Error(), "update failed")
	assert.NotContains(t, err.Error(), "local-skill1")
}

func TestRunSkillUpdate_AlreadyCurrent(t *testing.T) {
	homeDir := setupMockHome(t)

	destDir1 := filepath.Join(homeDir, ".agents", "skills", "current-skill")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "current-skill", OwnerRepo: "dummy/repo", SourceRevision: "sha-123"}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commit := skill.GitHubCommit{Sha: "sha-123"}
		require.NoError(t, json.NewEncoder(w).Encode(commit))
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "current-skill"})
	assert.NoError(t, err) // Already current, should not error
}

func TestRunSkillUpdate_MultipleFailures(t *testing.T) {
	homeDir := setupMockHome(t)

	destDir1 := filepath.Join(homeDir, ".agents", "skills", "bad1")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "bad1", OwnerRepo: "invalid/repo1"}))

	destDir2 := filepath.Join(homeDir, ".agents", "skills", "bad2")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "bad2", OwnerRepo: "invalid/repo2"}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bad1")
	assert.Contains(t, err.Error(), "bad2")
}
