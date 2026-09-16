package skill

import (
	"os"
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
		{"common", ".agents/skills"},
		{"copilot", ".github/copilot/skills"},
		{"cursor", ".cursor/skills"},
	}

	for _, tt := range tests {
		t.Run(tt.agent, func(t *testing.T) {
			target, err := ResolveTarget("user", tt.agent)
			assert.NoError(t, err)
			assert.Equal(t, tt.agent, target.Agent)
			assert.Equal(t, "user", target.Scope)
			assert.True(t, strings.HasSuffix(target.Path, tt.expected), "expected path to end with %s, got %s", tt.expected, target.Path)
			assert.True(t, strings.HasPrefix(target.Path, homeDir), "expected path to start with %s, got %s", homeDir, target.Path)
		})
	}
}

func TestResolveTarget_ProjectScope(t *testing.T) {
	// Simple test to ensure it creates a valid path in project scope
	target, err := ResolveTarget("project", "common")
	assert.NoError(t, err)
	assert.Equal(t, "common", target.Agent)
	assert.Equal(t, "project", target.Scope)
	assert.True(t, strings.HasSuffix(target.Path, ".agents/skills"))
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
