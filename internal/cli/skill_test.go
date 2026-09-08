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

func TestRunSkillUpdate_AllSuccess(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	os.MkdirAll(destDir, 0755)
	os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("ok"), 0644)

	skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "local-skill1",
		OriginalSource: "local",
	})

	err := RunSkillUpdate([]string{"--scope", "user", "local-skill1"})
	assert.NoError(t, err, "Updating a local skill should yield unsupported/local-only and not fail")
}

func TestRunSkillUpdate_InspectionFailure(t *testing.T) {
	err := RunSkillUpdate([]string{"--scope", "user", "non-existent-skill"})
	assert.Error(t, err, "Should fail inspection")
	assert.Contains(t, err.Error(), "skill 'non-existent-skill' not found")
}

func TestRunSkillUpdate_PartialFailure(t *testing.T) {
	homeDir := setupMockHome(t)

	// Create one valid local skill
	destDir1 := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	os.MkdirAll(destDir1, 0755)
	os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644)
	skill.SaveMetadata(destDir1, &skill.Metadata{Name: "local-skill1", OriginalSource: "local"})

	// Create one broken skill (missing metadata will fail inspect)
	// Actually, ListInstalledSkills only returns it if it can successfully parse the metadata.
	// So instead of --all, let's explicitly request a broken skill by name to force an inspection failure.
	destDir2 := filepath.Join(homeDir, ".agents", "skills", "broken-skill")
	os.MkdirAll(destDir2, 0755)
	os.WriteFile(filepath.Join(destDir2, ".rntocase-skill.json"), []byte("{bad json"), 0644)

	// Here we want to simulate an error in the RunSkillUpdate process.
	// If we provide the specific bad skill, it returns error early for single updates.
	// But the PR issue mentions `--all` behavior.
	// Let's mock a skill that parses but fails CheckUpdate or extract.
	skill.SaveMetadata(destDir2, &skill.Metadata{Name: "broken-skill", OwnerRepo: "invalid/repo"})

	// By making the API call fail, we can simulate an update failure.
	// But let's just make it simpler by asking for an explicit bad API call?
	// The problem is that ListInstalledSkills *hides* parse errors.
	// Since we want to test partial failure on `--all` let's have one valid and one network-failing skill.
	// We'll give it an invalid repo to ensure `CheckUpdate` fails.

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "broken-skill")
	// local-skill1 shouldn't be in the error output
	assert.NotContains(t, err.Error(), "local-skill1")
}

func TestRunSkillUpdate_AlreadyCurrent(t *testing.T) {
	homeDir := setupMockHome(t)

	destDir1 := filepath.Join(homeDir, ".agents", "skills", "current-skill")
	os.MkdirAll(destDir1, 0755)
	os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644)
	skill.SaveMetadata(destDir1, &skill.Metadata{Name: "current-skill", OwnerRepo: "dummy/repo", SourceRevision: "sha-123"})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commit := skill.GitHubCommit{Sha: "sha-123"}
		json.NewEncoder(w).Encode(commit)
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
	os.MkdirAll(destDir1, 0755)
	os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644)
	skill.SaveMetadata(destDir1, &skill.Metadata{Name: "bad1", OwnerRepo: "invalid/repo1"})

	destDir2 := filepath.Join(homeDir, ".agents", "skills", "bad2")
	os.MkdirAll(destDir2, 0755)
	os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("ok"), 0644)
	skill.SaveMetadata(destDir2, &skill.Metadata{Name: "bad2", OwnerRepo: "invalid/repo2"})

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bad1")
	assert.Contains(t, err.Error(), "bad2")
}
