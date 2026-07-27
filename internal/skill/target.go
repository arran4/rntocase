package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

// SupportedAgents defines the list of AI agents we know how to install skills for.
var SupportedAgents = []string{
	"common",
	"copilot",
	"cursor",
	"codex",
	"claude",
}

// Target defines the resolved installation path for a skill.
type Target struct {
	Agent string
	Scope string
	Path  string
}

// ResolveTarget finds the correct installation directory based on scope and agent.
// scope can be "user" or "project".
// agent can be empty (auto/common) or a specific agent name like "copilot", "cursor", "codex", "claude".
func ResolveTarget(scope string, agent string) (*Target, error) {
	if scope == "" {
		scope = "project" // Default to project scope
	}
	if agent == "" {
		agent = "common"
	}

	var basePath string
	if scope == "user" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		// In a real application, these paths would be well-documented or standardized
		// by the agents. Here we simulate a convention.
		switch agent {
		case "copilot":
			basePath = filepath.Join(homeDir, ".github", "copilot", "skills")
		case "cursor":
			basePath = filepath.Join(homeDir, ".cursor", "skills")
		case "codex":
			basePath = filepath.Join(homeDir, ".codex", "skills")
		case "claude":
			basePath = filepath.Join(homeDir, ".claude", "skills")
		case "common":
			basePath = filepath.Join(homeDir, ".agents", "skills")
		default:
			return nil, fmt.Errorf("unsupported agent for user scope: %s", agent)
		}
	} else if scope == "project" {
		// Find project root by looking for .git or fallback to current directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
		projectRoot := findProjectRoot(cwd)

		// Project scope usually uses a `.agents` or `.skills` folder in the repo root
		switch agent {
		case "copilot":
			basePath = filepath.Join(projectRoot, ".github", "copilot", "skills")
		case "cursor":
			basePath = filepath.Join(projectRoot, ".cursor", "skills")
		case "codex":
			basePath = filepath.Join(projectRoot, ".codex", "skills")
		case "claude":
			basePath = filepath.Join(projectRoot, ".claude", "skills")
		case "common":
			basePath = filepath.Join(projectRoot, ".agents", "skills")
		default:
			return nil, fmt.Errorf("unsupported agent for project scope: %s", agent)
		}
	} else {
		return nil, fmt.Errorf("invalid scope: %s. Must be 'user' or 'project'", scope)
	}

	return &Target{
		Agent: agent,
		Scope: scope,
		Path:  basePath,
	}, nil
}

// findProjectRoot attempts to find the root of the project (e.g., where .git is).
// If not found, it returns the original directory.
func findProjectRoot(dir string) string {
	currentDir := dir
	for {
		if _, err := os.Stat(filepath.Join(currentDir, ".git")); err == nil {
			return currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached root of filesystem
			break
		}
		currentDir = parentDir
	}
	return dir
}
