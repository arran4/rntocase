package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const MetadataFileName = ".rntocase-skill.json"

// Metadata represents the immutable source tracking information for an installed skill.
type Metadata struct {
	Name            string    `json:"name"`
	OriginalSource  string    `json:"original_source"`
	OwnerRepo       string    `json:"owner_repo,omitempty"` // For GitHub/remote sources
	PathWithin      string    `json:"path_within,omitempty"` // For subdirectories in a repo
	SourceRevision  string    `json:"source_revision,omitempty"`
	InstallTime     time.Time `json:"install_time"`
	InstallerApp    string    `json:"installer_app"`
	ContentDigest   string    `json:"content_digest,omitempty"`
}

// LoadMetadata reads the metadata file from the installed skill directory.
func LoadMetadata(skillDir string) (*Metadata, error) {
	metaPath := filepath.Join(skillDir, MetadataFileName)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("metadata not found in %s", skillDir)
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &meta, nil
}

// SaveMetadata writes the metadata to the installed skill directory.
func SaveMetadata(skillDir string, meta *Metadata) error {
	metaPath := filepath.Join(skillDir, MetadataFileName)

	// Create directory if it doesn't exist just to be safe
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}
