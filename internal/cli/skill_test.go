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
	err := RunSkill()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "requires a subcommand")
}

func TestRunSkill_UnknownSubcommand(t *testing.T) {
	err := RunSkill()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skill command requires a subcommand")
}

func TestRunSkillInstall_RequiresSource(t *testing.T) {
	err := RunSkillInstall("", "", false, "", "", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "remote source must be in owner/repo format")
}

func TestRunSkillUpdate_RequiresName(t *testing.T) {
	err := RunSkillUpdate("", "", false, false, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill update <name> or skill update --all")
}

func TestRunSkillRemove_RequiresName(t *testing.T) {
	err := RunSkillRemove("", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill remove <name>")
}

func TestRunSkillInspect_RequiresName(t *testing.T) {
	err := RunSkillInspect("", "", false, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usage: skill inspect <name>")
}

func TestRunSkillInstall_BasicLocal(t *testing.T) {
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

func TestRunSkillInstall_BundledOfficial(t *testing.T) {
	homeDir := setupMockHome(t)

	err := RunSkillInstall("user", "", false, "", "", "", "rntocase", "")
	require.NoError(t, err)

	destDir := filepath.Join(homeDir, ".agents", "skills", "rntocase")
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: rntocase")
}

func TestRunSkillUpdate_ExplicitForm(t *testing.T) {
	_ = setupMockHome(t)
	err := RunSkillInstall("user", "", false, "", "", "", "rntocase", "")
	require.NoError(t, err)

	err = RunSkillUpdate("user", "", false, false, "rntocase")
	require.NoError(t, err)
}

func TestRunSkillUpdate_LocalSkillError(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))

	require.NoError(t, skill.SaveMetadata(destDir, &skill.Metadata{
		Name:           "local-skill1",
		OriginalSource: "local",
	}))

	err := RunSkillUpdate("user", "", false, false, "local-skill1")
	assert.Error(t, err, "Updating a directly targeted local skill should preserve previous single-skill error semantics")
	assert.Contains(t, err.Error(), "is locally installed and cannot be updated automatically")
}

func TestRunSkillUpdate_InspectionFailure(t *testing.T) {
	err := RunSkillUpdate("user", "", false, false, "non-existent-skill")
	assert.Error(t, err, "Should fail inspection")
	assert.Contains(t, err.Error(), "skill 'non-existent-skill' not found")
}

func TestRunSkillUpdate_All_InspectionFailure(t *testing.T) {
	homeDir := setupMockHome(t)

	// Valid skill
	destDir1 := filepath.Join(homeDir, ".agents", "skills", "local-skill1")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "local-skill1", OriginalSource: "local"}))

	// Second skill that ListInstalledSkills can find, but InspectSkill will fail on.
	// This happens if the directory name doesn't match the inner skill name.
	destDir2 := filepath.Join(homeDir, ".agents", "skills", "mismatched-dir")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "inner-name", OriginalSource: "local"}))

	err := RunSkillUpdate("user", "", false, true, "")
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
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "a-broken", OwnerRepo: "invalid/repo"}))

	// Lexically-later valid remote skill
	destDir2 := filepath.Join(homeDir, ".agents", "skills", "z-current")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
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

	err := RunSkillUpdate("user", "", false, true, "")
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
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "current-skill", OwnerRepo: "dummy/repo", SourceRevision: "sha-123"}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		commit := skill.GitHubCommit{Sha: "sha-123"}
		require.NoError(t, json.NewEncoder(w).Encode(commit))
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate("user", "", false, false, "current-skill")
	assert.NoError(t, err) // Already current, should not error
}

func TestRunSkillUpdate_MultipleFailures(t *testing.T) {
	homeDir := setupMockHome(t)

	destDir1 := filepath.Join(homeDir, ".agents", "skills", "bad1")
	require.NoError(t, os.MkdirAll(destDir1, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir1, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir1, &skill.Metadata{Name: "bad1", OwnerRepo: "invalid/repo1"}))

	destDir2 := filepath.Join(homeDir, ".agents", "skills", "bad2")
	require.NoError(t, os.MkdirAll(destDir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir2, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(destDir2, &skill.Metadata{Name: "bad2", OwnerRepo: "invalid/repo2"}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate("user", "", false, true, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bad1")
	assert.Contains(t, err.Error(), "bad2")
}

func TestRunSkillInstall_WithFlags(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "custom-skill-name")

	// Create a mock local source with a SKILL.md
	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))

	err := RunSkillInstall("user", "", false, "", "", "custom-skill-name", sourceDir, "")
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "---\nname: custom-skill-name\ndescription: test skill\n---", string(content))
}

func TestRunSkillInstall_WithFlags_PathTraversal(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall("user", "", false, "", "", "../escaped", sourceDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")
}
func TestRunSkillInstall_DotName(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall("user", "", false, "", "", ".", sourceDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")
}

func TestRunSkillInstall_PathTraversal(t *testing.T) {
	_ = setupMockHome(t)
	sourceDir := t.TempDir()

	err := RunSkillInstall("user", "", false, "", "", "..", sourceDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")
}

func TestRunSkillUpdate_PinnedRef(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "pinned-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))

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

	err := RunSkillUpdate("user", "", false, false, "pinned-skill")
	assert.NoError(t, err) // Already current, should not error
}

// RouteAssertingTransport is a mock http.RoundTripper that returns a specific status code
// while enforcing exact HTTP method and URL paths.
type RouteAssertingTransport struct {
	ExpectedMethod string
	ExpectedURL    string
	StatusCode     int
	t              *testing.T
}

func (t *RouteAssertingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != t.ExpectedMethod {
		t.t.Errorf("expected HTTP method %s, got %s", t.ExpectedMethod, req.Method)
	}
	actualURL := req.URL.String()
	if actualURL != t.ExpectedURL {
		t.t.Errorf("expected URL %s, got %s", t.ExpectedURL, actualURL)
	}

	return &http.Response{
		StatusCode: t.StatusCode,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}, nil
}

func TestRunSkillInstall_RefNotFound(t *testing.T) {
	expectedAPIURL := skill.GitHubAPIURL + "/repos/arran4/mock-rntocase-notfound/commits/does-not-exist"

	originalClient := skill.HTTPClient
	skill.HTTPClient = &http.Client{
		Transport: &RouteAssertingTransport{
			ExpectedMethod: "GET",
			ExpectedURL:    expectedAPIURL,
			StatusCode:     http.StatusNotFound,
			t:              t,
		},
	}
	defer func() { skill.HTTPClient = originalClient }()

	_ = setupMockHome(t)
	err := RunSkillInstall("user", "", false, "does-not-exist", "", "", "arran4/mock-rntocase-notfound", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get repository metadata: HTTP 404")
}

func TestRunSkillUpdate_TrackingRef(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "tracking-skill")
	require.NoError(t, os.MkdirAll(destDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("---\nname: custom-skill-name\ndescription: test skill\n---"), 0644))

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
			hdr := &tar.Header{Name: "repo-sha/SKILL.md", Mode: 0600, Size: 52}
			if err := tw.WriteHeader(hdr); err != nil {
				panic(err)
			}
			if _, err := tw.Write([]byte("---\nname: tracking-skill\ndescription: test skill\n---")); err != nil {
				panic(err)
			}
			if err := tw.Close(); err != nil {
				panic(err)
			}
			if err := gw.Close(); err != nil {
				panic(err)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillUpdate("user", "", false, false, "tracking-skill")
	assert.NoError(t, err)

	meta, err := skill.LoadMetadata(destDir)
	assert.NoError(t, err)
	assert.Equal(t, "new-sha", meta.SourceRevision)
}

func TestRunSkillInstall_NameValidation(t *testing.T) {
	homeDir := setupMockHome(t)
	destDir := filepath.Join(homeDir, ".agents", "skills", "good-name")

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: good-name\ndescription: desc\n---"), 0644))

	// Valid match
	err := RunSkillInstall("user", "", false, "", "", "good-name", sourceDir, "")
	assert.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))

	// Valid manifest but user requests mismatch name via --name
	sourceDirMismatch := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDirMismatch, "SKILL.md"), []byte("---\nname: other-name\ndescription: desc\n---"), 0644))
	err = RunSkillInstall("user", "", false, "", "", "mismatch-name", sourceDirMismatch, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "skill name in manifest ('other-name') does not match installed directory name ('mismatch-name')")

	// Invalid name via --name
	err = RunSkillInstall("user", "", false, "", "", "Bad-Name", sourceDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")

	// Invalid legacy positional name
	err = RunSkillInstall("user", "", false, "", "", "", sourceDir, "Bad-Name")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")

	// Invalid source inference
	err = RunSkillInstall("user", "", false, "", "", "", filepath.Join(sourceDir, "Bad-Name-Dir"), "")
	assert.Error(t, err)
	// It fails path classification or validation, but we expect validation if it's local
}

func TestRunSkillUpdate_AllCrossAgentRegression(t *testing.T) {
	homeDir := setupMockHome(t)
	expectedDir := t.TempDir()
	require.NoError(t, skill.ExtractEmbeddedSkill("rntocase", expectedDir))
	expectedSkill, err := os.ReadFile(filepath.Join(expectedDir, "SKILL.md"))
	require.NoError(t, err)
	expectedDigest, err := skill.ComputeDirectoryDigest(expectedDir)
	require.NoError(t, err)

	// Setup same skill name under both copilot and cursor roots
	copilotDir := filepath.Join(homeDir, ".copilot", "skills", "rntocase")
	require.NoError(t, os.MkdirAll(copilotDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(copilotDir, "SKILL.md"), []byte("---\nname: rntocase\ndescription: desc\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(copilotDir, &skill.Metadata{
		Name:           "rntocase",
		OriginalSource: "official", // use official to avoid network mock
		ContentDigest:  "old-digest-copilot",
	}))

	cursorDir := filepath.Join(homeDir, ".cursor", "skills", "rntocase")
	require.NoError(t, os.MkdirAll(cursorDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(cursorDir, "SKILL.md"), []byte("---\nname: rntocase\ndescription: desc\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(cursorDir, &skill.Metadata{
		Name:           "rntocase",
		OriginalSource: "official",
		ContentDigest:  "old-digest-cursor",
	}))
	require.NoError(t, os.WriteFile(filepath.Join(cursorDir, "cursor-sentinel.txt"), []byte("cursor must remain untouched"), 0644))

	// Capture the entire non-target installation before updating Copilot.
	cursorSkillBefore, err := os.ReadFile(filepath.Join(cursorDir, "SKILL.md"))
	require.NoError(t, err)
	cursorMetadataBefore, err := os.ReadFile(filepath.Join(cursorDir, skill.MetadataFileName))
	require.NoError(t, err)
	cursorSentinelBefore, err := os.ReadFile(filepath.Join(cursorDir, "cursor-sentinel.txt"))
	require.NoError(t, err)

	// Run update --all for copilot agent only
	err = RunSkillUpdate("user", "copilot", true, true, "")
	assert.NoError(t, err)

	// Verify copilot was updated (ContentDigest changed)
	copilotSkill, err := os.ReadFile(filepath.Join(copilotDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, expectedSkill, copilotSkill)
	copilotMeta, err := skill.LoadMetadata(copilotDir)
	require.NoError(t, err)
	assert.NotEqual(t, "old-digest-copilot", copilotMeta.ContentDigest)
	assert.Equal(t, "rntocase", copilotMeta.Name)
	assert.Equal(t, "official", copilotMeta.OriginalSource)
	assert.Equal(t, expectedDigest, copilotMeta.ContentDigest)

	// Verify the Cursor installation was not updated in any way.
	cursorSkillAfter, err := os.ReadFile(filepath.Join(cursorDir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, cursorSkillBefore, cursorSkillAfter)
	cursorMetadataAfter, err := os.ReadFile(filepath.Join(cursorDir, skill.MetadataFileName))
	require.NoError(t, err)
	assert.Equal(t, cursorMetadataBefore, cursorMetadataAfter)
	cursorSentinelAfter, err := os.ReadFile(filepath.Join(cursorDir, "cursor-sentinel.txt"))
	require.NoError(t, err)
	assert.Equal(t, cursorSentinelBefore, cursorSentinelAfter)
}

func TestResolveSkillPath_TraversalRegression(t *testing.T) {
	homeDir := setupMockHome(t)
	// Create an external directory outside the agents path, simulating another place with a skill
	externalDir := filepath.Join(homeDir, "external-malicious")
	require.NoError(t, os.MkdirAll(externalDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(externalDir, "SKILL.md"), []byte("---\nname: malicious\ndescription: desc\n---"), 0644))
	require.NoError(t, skill.SaveMetadata(externalDir, &skill.Metadata{Name: "malicious"}))
	require.NoError(t, os.WriteFile(filepath.Join(externalDir, "external-sentinel.txt"), []byte("do not modify"), 0644))

	externalSkillBefore, err := os.ReadFile(filepath.Join(externalDir, "SKILL.md"))
	require.NoError(t, err)
	externalMetadataBefore, err := os.ReadFile(filepath.Join(externalDir, skill.MetadataFileName))
	require.NoError(t, err)
	externalSentinelBefore, err := os.ReadFile(filepath.Join(externalDir, "external-sentinel.txt"))
	require.NoError(t, err)
	assertExternalUnchanged := func() {
		externalSkillAfter, readErr := os.ReadFile(filepath.Join(externalDir, "SKILL.md"))
		require.NoError(t, readErr)
		assert.Equal(t, externalSkillBefore, externalSkillAfter)
		externalMetadataAfter, readErr := os.ReadFile(filepath.Join(externalDir, skill.MetadataFileName))
		require.NoError(t, readErr)
		assert.Equal(t, externalMetadataBefore, externalMetadataAfter)
		externalSentinelAfter, readErr := os.ReadFile(filepath.Join(externalDir, "external-sentinel.txt"))
		require.NoError(t, readErr)
		assert.Equal(t, externalSentinelBefore, externalSentinelAfter)
	}

	// Attempt RemoveSkill with traversal: we are at `~/.agents/skills/`, so `../../external-malicious` should escape.
	err = RunSkillRemove("user", "common", "../../external-malicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name must be a single component without separators")

	// Ensure the external directory still exists
	_, err = os.Stat(externalDir)
	assert.NoError(t, err, "external directory should not be removed by traversal attack")
	assertExternalUnchanged()

	err = RunSkillInspect("user", "common", false, "../../external-malicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name must be a single component without separators")
	assertExternalUnchanged()

	err = RunSkillUpdate("user", "common", false, false, "../../external-malicious")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name must be a single component without separators")
	assertExternalUnchanged()
}
