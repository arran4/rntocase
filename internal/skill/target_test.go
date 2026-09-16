package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveTarget_UserScope(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		agent    string
		expected string
	}{
		{"copilot", filepath.Join(homeDir, ".copilot", "skills")},
		{"cursor", filepath.Join(homeDir, ".cursor", "skills")},
		{"codex", filepath.Join(homeDir, ".codex", "skills")},
		{"claude", filepath.Join(homeDir, ".claude", "skills")},
		{"common", filepath.Join(homeDir, ".agents", "skills")},
		{"", filepath.Join(homeDir, ".agents", "skills")},
	}

	for _, tt := range tests {
		t.Run(tt.agent, func(t *testing.T) {
			target, err := ResolveTarget("user", tt.agent)
			assert.NoError(t, err)
			expectedAgent := tt.agent
			if expectedAgent == "" {
				expectedAgent = "common"
			}
			assert.Equal(t, expectedAgent, target.Agent)
			assert.Equal(t, "user", target.Scope)
			assert.Equal(t, tt.expected, target.Path)
			assert.True(t, strings.HasPrefix(target.Path, homeDir), "expected path to start with %s, got %s", homeDir, target.Path)
		})
	}
}

func TestResolveTarget_ProjectScope(t *testing.T) {
	tests := []struct {
		agent    string
		expected string
	}{
		{"copilot", ".github/skills"},
		{"cursor", ".cursor/skills"},
		{"codex", ".codex/skills"},
		{"claude", ".claude/skills"},
		{"common", ".agents/skills"},
		{"", ".agents/skills"},
	}

	for _, tt := range tests {
		t.Run(tt.agent, func(t *testing.T) {
			target, err := ResolveTarget("project", tt.agent)
			assert.NoError(t, err)
			expectedAgent := tt.agent
			if expectedAgent == "" {
				expectedAgent = "common"
			}
			assert.Equal(t, expectedAgent, target.Agent)
			assert.Equal(t, "project", target.Scope)
			expectedSuffix := filepath.FromSlash(tt.expected)
			assert.True(t, strings.HasSuffix(target.Path, expectedSuffix), "expected path to end with %s, got %s", expectedSuffix, target.Path)
		})
	}
}

func TestResolveTarget_InvalidScope(t *testing.T) {
	_, err := ResolveTarget("invalid", "common")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid scope")
}

func TestResolveTarget_InvalidAgent(t *testing.T) {
	_, err := ResolveTarget("user", "invalid_agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported agent")
}

func TestResolveSkillPath(t *testing.T) {
	target := &Target{
		Path: "/base/path",
	}

	tests := []struct {
		name      string
		skillName string
		wantErr   bool
	}{
		{"valid name", "my-skill", false},
		{"empty name", "", true},
		{"traversal", "../escaped", true},
		{"nested valid", "nested/skill", false},
		{"absolute path", "/absolute/skill", true},
		{"traversal tricky", "my-skill/../../escaped", true},
		{"same as root", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveSkillPath(target, tt.skillName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveSkillPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
