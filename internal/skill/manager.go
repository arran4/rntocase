package skill

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	// GitHubAPIURL is the base URL for GitHub API requests.
	GitHubAPIURL = "https://api.github.com"
	// HTTPClient is the standard client for skill network operations, providing bounded timeouts.
	HTTPClient = &http.Client{Timeout: 15 * time.Second}
)

// InstalledSkillInfo holds metadata and the agent it was found under.
type InstalledSkillInfo struct {
	Meta  *Metadata
	Agent string
}

// ListInstalledSkills searches the known agent paths for installed skills and returns their info.
func ListInstalledSkills(scope string) ([]*InstalledSkillInfo, error) {
	var results []*InstalledSkillInfo

	for _, agent := range SupportedAgents {
		target, err := ResolveTarget(scope, agent)
		if err != nil {
			continue // Skip if not supported or error
		}

		err = filepath.Walk(target.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip errors like permission denied
			}

			if !info.IsDir() {
				if info.Name() == MetadataFileName {
					skillDir := filepath.Dir(path)
					meta, err := LoadMetadata(skillDir)
					if err == nil {
						results = append(results, &InstalledSkillInfo{
							Meta:  meta,
							Agent: agent,
						})
					}
				}
			}
			return nil
		})

		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read agent directory %s: %w", target.Path, err)
		}
	}

	return results, nil
}

// InspectSkill looks up a specific skill by name and scope and returns its details.
func InspectSkill(name string, scope string, agent string) (*Metadata, string, error) {
	target, err := ResolveTarget(scope, agent)
	if err != nil {
		return nil, "", err
	}

	skillDir := filepath.Join(target.Path, name)
	meta, err := LoadMetadata(skillDir)
	if err != nil {
		return nil, "", fmt.Errorf("skill '%s' not found or invalid: %w", name, err)
	}

	return meta, skillDir, nil
}

// RemoveSkill removes a skill directory if it exists and contains valid metadata.
func RemoveSkill(name string, scope string, agent string) error {
	target, err := ResolveTarget(scope, agent)
	if err != nil {
		return err
	}

	skillDir := filepath.Join(target.Path, name)

	// Safety check: only remove if it has our metadata file
	if _, err := LoadMetadata(skillDir); err != nil {
		return fmt.Errorf("skill '%s' not found or is not managed by this tool: %w", name, err)
	}

	if err := os.RemoveAll(skillDir); err != nil {
		return fmt.Errorf("failed to remove skill directory: %w", err)
	}

	return nil
}

// CheckUpdate checks if an update is available for the given skill metadata.
func CheckUpdate(meta *Metadata) (bool, string, error) {
	if meta.OwnerRepo == "" {
		return false, "", fmt.Errorf("local skills cannot be updated automatically")
	}

	// If the requested ref is a 40-character hex string (an exact SHA pin), it never updates.
	if len(meta.RequestedRef) == 40 {
		isHex := true
		for _, c := range meta.RequestedRef {
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				isHex = false
				break
			}
		}
		if isHex {
			return false, meta.SourceRevision, nil
		}
	}

	ref := meta.RequestedRef
	if ref == "" {
		ref = "HEAD"
	}

	apiURL := fmt.Sprintf("%s/repos/%s/commits/%s", GitHubAPIURL, meta.OwnerRepo, ref)
	// Optionally add path parameter if it's a subfolder? We just check the whole repo HEAD here
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, "", err
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch upstream status: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("failed to fetch upstream status: HTTP %d", resp.StatusCode)
	}

	var commit GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return false, "", err
	}

	return commit.Sha != meta.SourceRevision, commit.Sha, nil
}
