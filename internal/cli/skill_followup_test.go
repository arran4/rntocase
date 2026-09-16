package cli

import (
	"archive/tar"
	"compress/gzip"
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

func TestRunSkillInstallRejectsInvalidInferredLocalName(t *testing.T) {
	_ = setupMockHome(t)

	sourceDir := filepath.Join(t.TempDir(), "Bad-Name-Dir")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "SKILL.md"), []byte("---\nname: valid-name\ndescription: desc\n---\n"), 0644))

	err := RunSkillInstall("user", "", false, "", "", "", sourceDir, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid skill name: name must contain only lowercase ASCII letters")
}

func TestRunSkillInstallLegacySecondPositionalRepositoryPath(t *testing.T) {
	homeDir := setupMockHome(t)
	const manifest = "---\nname: example\ndescription: desc\n---\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/owner/repo/commits/HEAD":
			require.NoError(t, json.NewEncoder(w).Encode(skill.GitHubCommit{Sha: "sha-123"}))
		case "/repos/owner/repo/tarball/sha-123":
			w.WriteHeader(http.StatusOK)
			gw := gzip.NewWriter(w)
			tw := tar.NewWriter(gw)
			hdr := &tar.Header{
				Name: "repo-sha/skills/example/SKILL.md",
				Mode: 0644,
				Size: int64(len(manifest)),
			}
			require.NoError(t, tw.WriteHeader(hdr))
			_, err := tw.Write([]byte(manifest))
			require.NoError(t, err)
			require.NoError(t, tw.Close())
			require.NoError(t, gw.Close())
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	originalAPIURL := skill.GitHubAPIURL
	skill.GitHubAPIURL = ts.URL
	defer func() { skill.GitHubAPIURL = originalAPIURL }()

	err := RunSkillInstall("user", "", false, "", "", "", "owner/repo", "skills/example")
	require.NoError(t, err)

	destDir := filepath.Join(homeDir, ".agents", "skills", "example")
	assert.FileExists(t, filepath.Join(destDir, "SKILL.md"))
	meta, err := skill.LoadMetadata(destDir)
	require.NoError(t, err)
	assert.Equal(t, "skills/example", meta.PathWithin)
	assert.Equal(t, "owner/repo", meta.OwnerRepo)
}
