package cli

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	assert.Contains(t, err.Error(), "usage: skill install [--replace] [--ref <ref>] [--path <path>] [--name <name>] <source>")
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

	// Lexically-first failing remote skill
	destDir1 := filepath.Join(homeDir, ".agents", "skills", "a-broken")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "a-broken", OwnerRepo: "invalid/repo"}))

	// Lexically-later valid remote skill
	destDir2 := filepath.Join(homeDir, ".agents", "skills", "z-current")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("ok"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "z-current", OwnerRepo: "valid/repo", SourceRevision: "sha-123"}))

	var validRepoCalled bool

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/invalid/repo/commits/HEAD" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.Path == "/repos/valid/repo/commits/HEAD" {
			validRepoCalled = true
			commit := skill.GitHubCommit{Sha: "sha-123"}
			require.NoError(t, json.NewEncoder(w).Encode(commit))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "--all"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "a-broken")
	assert.Contains(t, err.Error(), "update failed")
	assert.NotContains(t, err.Error(), "z-current")

	// Crucial assertion: Verify the loop continued and the valid repo was still checked!
	assert.True(t, validRepoCalled, "The valid remote skill should have been processed even though a prior skill failed")
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

func TestRunSkillInstall_WithFlags(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "custom-skill-name")

	// Create a mock local source with a SKILL.md
	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("ok"), 0644))

	err := RunSkillInstall([]string{"--scope", "user", "--name", "custom-skill-name", sourceDir})
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(content))
}

func TestRunSkillInstall_WithFlags_PathTraversal(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall([]string{"--scope", "user", "--name", "../escaped", sourceDir})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: ../escaped")
}
func TestRunSkillInstall_DotName(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall([]string{"--scope", "user", "--name", ".", sourceDir})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: .")
}

func TestRunSkillInstall_PathTraversal(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall([]string{"--scope", "user", "--name", "..", sourceDir})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: ..")
}

func TestRunSkillUpdate_PinnedRef(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "pinned-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("ok"), 0644))

	require.NoError(t, skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "pinned-skill",
		OriginalSource: "dummy/repo",
		OwnerRepo:      "dummy/repo",
		SourceRevision: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
		RequestedRef:   "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
	}))

	// Start a dummy server that panics if called, to prove CheckUpdate exits early
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Update check should have been skipped for exact SHA pin")
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "pinned-skill"})
	assert.NoError(t, err) // Already current, should not error
}

func TestRunSkillInstall_RefNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	_ = setupMockHome(t)
	err := RunSkillInstall([]string{"--scope", "user", "--ref", "does-not-exist", "dummy/repo"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get repository metadata: HTTP 404")
}

func TestRunSkillUpdate_TrackingRef(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "tracking-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("ok"), 0644))

	require.NoError(t, skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "tracking-skill",
		OriginalSource: "dummy/repo",
		OwnerRepo:      "dummy/repo",
		SourceRevision: "old-sha",
		RequestedRef:   "main",
	}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "commits/main") {
			commit := skill.GitHubCommit{Sha: "new-sha"}
			_ = json.NewEncoder(w).Encode(commit)
			return
		}
		if strings.HasSuffix(r.URL.Path, "tarball/new-sha") {
			w.WriteHeader(http.StatusOK)
			// Return dummy tarball with SKILL.md
			gw := gzip.NewWriter(w)
			tw := tar.NewWriter(gw)
			hdr := &tar.Header{Name: "repo-sha/SKILL.md", Mode: 0600, Size: 2}
			tw.WriteHeader(hdr)
			tw.Write([]byte("ok"))
			tw.Close()
			gw.Close()
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate([]string{"--scope", "user", "tracking-skill"})
	assert.NoError(t, err)

	meta, err := skill.LoadMetadata(destDir)
	assert.NoError(t, err)
	assert.Equal(t, "new-sha", meta.SourceRevision)
}
