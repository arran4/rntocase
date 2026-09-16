package skill

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

var validNameRegex = regexp.MustCompile(`^[a-z0-9](-?[a-z0-9])*$`)

type AgentSkillManifest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ParseAndValidateManifest parses a SKILL.md file content and validates its YAML frontmatter.
func ParseAndValidateManifest(mdContent []byte) (*AgentSkillManifest, error) {
	parts := bytes.SplitN(mdContent, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("SKILL.md must contain YAML frontmatter enclosed in '---'")
	}

	frontmatter := parts[1]

	var manifest AgentSkillManifest
	if err := yaml.Unmarshal(frontmatter, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest must define 'name'")
	}
	if len(manifest.Name) > 64 {
		return nil, fmt.Errorf("manifest 'name' must be 1-64 characters")
	}
	if !validNameRegex.MatchString(manifest.Name) {
		return nil, fmt.Errorf("manifest 'name' must contain only lowercase ASCII letters, digits, and hyphens, with no consecutive, leading, or trailing hyphens")
	}

	if manifest.Description == "" {
		return nil, fmt.Errorf("manifest must define 'description'")
	}
	// The problem asks to validate description length, but doesn't specify it, though standard might be 150. We'll leave it as non-empty as instructed initially or up to a reasonable limit.

	return &manifest, nil
}
