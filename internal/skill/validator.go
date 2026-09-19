package skill

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var validNameRegex = regexp.MustCompile(`^[a-z0-9](-?[a-z0-9])*$`)

type AgentSkillManifest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ParseAndValidateManifest parses a SKILL.md file content and validates its YAML frontmatter.
func ParseAndValidateManifest(mdContent []byte) (*AgentSkillManifest, error) {
	normalized := bytes.ReplaceAll(mdContent, []byte("\r\n"), []byte("\n"))

	if !bytes.HasPrefix(normalized, []byte("---\n")) {
		return nil, fmt.Errorf("SKILL.md must start with exactly '---' on the first line")
	}

	endIdx := bytes.Index(normalized[4:], []byte("\n---\n"))
	var frontmatter []byte
	if endIdx == -1 {
		if string(normalized) == "---\n---" || string(normalized) == "---\n---\n" {
			frontmatter = []byte{}
		} else if bytes.HasSuffix(normalized, []byte("\n---")) {
			frontmatter = normalized[4 : len(normalized)-4]
		} else {
			return nil, fmt.Errorf("SKILL.md missing closing '---' on a standalone line")
		}
	} else {
		frontmatter = normalized[4 : 4+endIdx]
	}

	var manifest AgentSkillManifest
	if err := yaml.Unmarshal(frontmatter, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	if err := ValidateSkillName(manifest.Name); err != nil {
		return nil, fmt.Errorf("manifest 'name' validation failed: %w", err)
	}

	manifest.Description = strings.TrimSpace(manifest.Description)
	descLen := len([]rune(manifest.Description))
	if descLen < 1 || descLen > 1024 {
		return nil, fmt.Errorf("manifest 'description' must be 1-1024 characters")
	}

	return &manifest, nil
}

// ValidateSkillName explicitly validates an Agent Skills name
func ValidateSkillName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("name must be 1-64 characters")
	}
	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("name must contain only lowercase ASCII letters, digits, and hyphens, with no consecutive, leading, or trailing hyphens")
	}
	return nil
}
