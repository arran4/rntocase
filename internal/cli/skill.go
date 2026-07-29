package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arran4/rntocase/internal/skill"
)

// RunSkill is a subcommand `rntocase skill` -- Manage AI agent skills for this CLI
func RunSkill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("skill command requires a subcommand (install, update, remove, list, inspect)")
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "install":
		return RunSkillInstall(subArgs)
	case "update":
		return RunSkillUpdate(subArgs)
	case "remove":
		return RunSkillRemove(subArgs)
	case "list":
		return RunSkillList(subArgs)
	case "inspect":
		return RunSkillInspect(subArgs)
	default:
		return fmt.Errorf("unknown skill subcommand: %s", subcommand)
	}
}

// RunSkillInstall is a subcommand `rntocase skill install` -- Install a skill
func RunSkillInstall(args []string) error {
	fs := flag.NewFlagSet("skill install", flag.ExitOnError)
	scope := fs.String("scope", "project", "Installation scope: user or project")
	agent := fs.String("agent", "common", "Target agent: common, copilot, cursor, codex, claude")
	_ = fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		return fmt.Errorf("usage: skill install <source> [skill-name-or-path]")
	}

	source := positionalArgs[0]
	nameOrPath := ""
	if len(positionalArgs) > 1 {
		nameOrPath = positionalArgs[1]
	}

	target, err := skill.ResolveTarget(*scope, *agent)
	if err != nil {
		return err
	}

	isLocal := strings.HasPrefix(source, ".") || strings.HasPrefix(source, "/")
	isOfficial := source == "official" || source == "rntocase"

	// Determine final skill name
	skillName := nameOrPath
	if skillName == "" {
		if isLocal {
			skillName = filepath.Base(source)
		} else if isOfficial {
			skillName = "rntocase"
		} else {
			parts := strings.Split(source, "/")
			if len(parts) >= 2 {
				skillName = parts[1] // fallback to repo name
			} else {
				return fmt.Errorf("could not determine skill name automatically, please provide it")
			}
		}
	}

	destDir := filepath.Join(target.Path, skillName)

	// Ensure agent skills directory exists
	if err := os.MkdirAll(target.Path, 0755); err != nil {
		return fmt.Errorf("failed to create agent directory %s: %w", target.Path, err)
	}

	fmt.Printf("Installing skill '%s' to %s...\n", skillName, destDir)

	meta := &skill.Metadata{
		Name:           skillName,
		OriginalSource: source,
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
	}

	if isLocal {
		if err := skill.CopyLocalDirectory(source, destDir); err != nil {
			return fmt.Errorf("failed to copy local skill: %w", err)
		}
	} else if isOfficial {
		if err := skill.ExtractEmbeddedSkill("rntocase", destDir); err != nil {
			return fmt.Errorf("failed to install official skill: %w", err)
		}
	} else {
		// Remote Github
		ownerRepo := source
		if !strings.Contains(ownerRepo, "/") {
			return fmt.Errorf("remote source must be in owner/repo format (e.g. arran4/rntocase) or 'official'")
		}

		pathWithin := ""
		if nameOrPath != "" && strings.Contains(nameOrPath, "/") {
			pathWithin = nameOrPath
		}

		tarPath, sha, err := skill.DownloadGitHubRepository(ownerRepo)
		if err != nil {
			return fmt.Errorf("failed to download skill: %w", err)
		}
		defer func() { _ = os.Remove(tarPath) }()

		meta.OwnerRepo = ownerRepo
		meta.SourceRevision = sha
		meta.PathWithin = pathWithin

		if err := skill.ExtractTarGz(tarPath, destDir, pathWithin); err != nil {
			return fmt.Errorf("failed to extract skill: %w", err)
		}
	}

	// Validate SKILL.md
	if _, err := os.Stat(filepath.Join(destDir, "SKILL.md")); os.IsNotExist(err) {
		// Cleanup if invalid
		_ = os.RemoveAll(destDir)
		return fmt.Errorf("installation failed: skill must contain a SKILL.md file")
	}

	digest, err := skill.ComputeDirectoryDigest(destDir)
	if err == nil {
		meta.ContentDigest = digest
	}

	if err := skill.SaveMetadata(destDir, meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Println("Skill installed successfully.")
	return nil
}

// RunSkillUpdate is a subcommand `rntocase skill update` -- Update a skill
func RunSkillUpdate(args []string) error {
	fs := flag.NewFlagSet("skill update", flag.ExitOnError)
	scope := fs.String("scope", "project", "Installation scope")
	agent := fs.String("agent", "common", "Target agent")
	force := fs.Bool("force", false, "Force update and overwrite local changes")
	all := fs.Bool("all", false, "Update all installed skills in the given scope")
	_ = fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 && !*all {
		return fmt.Errorf("usage: skill update <name> or skill update --all")
	}

	var skillsToUpdate []struct {
		Name string
		Meta *skill.Metadata
		Dir  string
	}

	if *all {
		installed, err := skill.ListInstalledSkills(*scope)
		if err != nil {
			return err
		}
		for _, info := range installed {
			// Find the actual path using the agent it was found under
			_, dir, err := skill.InspectSkill(info.Meta.Name, *scope, info.Agent)
			if err == nil {
				skillsToUpdate = append(skillsToUpdate, struct {
					Name string
					Meta *skill.Metadata
					Dir  string
				}{info.Meta.Name, info.Meta, dir})
			}
		}
	} else {
		name := positionalArgs[0]
		meta, destDir, err := skill.InspectSkill(name, *scope, *agent)
		if err != nil {
			return err
		}
		skillsToUpdate = append(skillsToUpdate, struct {
			Name string
			Meta *skill.Metadata
			Dir  string
		}{name, meta, destDir})
	}

	if len(skillsToUpdate) == 0 {
		fmt.Println("No skills found to update.")
		return nil
	}

	for _, s := range skillsToUpdate {
		if err := updateSingleSkill(s.Name, s.Meta, s.Dir, *force); err != nil {
			fmt.Printf("Failed to update '%s': %v\n", s.Name, err)
		}
	}

	return nil
}

func updateSingleSkill(name string, meta *skill.Metadata, destDir string, force bool) error {
	if meta.OriginalSource == "official" || meta.OriginalSource == "rntocase" {
		fmt.Printf("Checking for updates for official skill '%s'...\n", name)

		// If there is an update to the binary it could bring a new embedded skill
		// Since we don't have a network version for embedded, we just forcefully re-extract it
		// Check for local mods first
		currentDigest, err := skill.ComputeDirectoryDigest(destDir)
		if err == nil && meta.ContentDigest != "" && currentDigest != meta.ContentDigest && !force {
			fmt.Printf("Skill '%s' has local modifications. Use --force to replace.\n", name)
			return nil
		}

		_ = os.RemoveAll(destDir)
		if err := skill.ExtractEmbeddedSkill("rntocase", destDir); err != nil {
			return err
		}

		digest, _ := skill.ComputeDirectoryDigest(destDir)
		meta.ContentDigest = digest
		meta.InstallTime = time.Now()
		if err := skill.SaveMetadata(destDir, meta); err != nil {
			return err
		}

		fmt.Printf("Skill '%s' updated from embedded source.\n", name)
		return nil
	}

	if meta.OwnerRepo == "" {
		return fmt.Errorf("skill '%s' is locally installed and cannot be updated automatically", name)
	}

	fmt.Printf("Checking for updates for '%s'...\n", name)
	hasUpdate, shaUpdate, err := skill.CheckUpdate(meta)
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Println("Skill is already up to date.")
		return nil
	}

	// Check for local modifications
	currentDigest, err := skill.ComputeDirectoryDigest(destDir)
	if err == nil && meta.ContentDigest != "" && currentDigest != meta.ContentDigest && !force {
		fmt.Printf("Skill '%s' has local modifications. Use --force to replace.\n", name)
		return nil
	}

	fmt.Println("Updating skill...")

	// Re-download first to be more atomic
	tarPath, shaDownload, err := skill.DownloadGitHubRepository(meta.OwnerRepo)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer func() { _ = os.Remove(tarPath) }()

	// Remove old
	_ = os.RemoveAll(destDir)

	if err := skill.ExtractTarGz(tarPath, destDir, meta.PathWithin); err != nil {
		return fmt.Errorf("failed to extract updated skill: %w", err)
	}

	meta.SourceRevision = shaDownload
	meta.InstallTime = time.Now()

	digest, err := skill.ComputeDirectoryDigest(destDir)
	if err == nil {
		meta.ContentDigest = digest
	}

	if err := skill.SaveMetadata(destDir, meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	shortSha := shaUpdate
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}

	fmt.Printf("Skill updated successfully to revision %s.\n", shortSha)
	return nil
}

// RunSkillRemove is a subcommand `rntocase skill remove` -- Remove a skill
func RunSkillRemove(args []string) error {
	fs := flag.NewFlagSet("skill remove", flag.ExitOnError)
	scope := fs.String("scope", "project", "Installation scope")
	agent := fs.String("agent", "common", "Target agent")
	_ = fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		return fmt.Errorf("usage: skill remove <name>")
	}

	name := positionalArgs[0]

	if err := skill.RemoveSkill(name, *scope, *agent); err != nil {
		return err
	}

	fmt.Printf("Skill '%s' removed successfully.\n", name)
	return nil
}

// RunSkillList is a subcommand `rntocase skill list` -- List skills
func RunSkillList(args []string) error {
	fs := flag.NewFlagSet("skill list", flag.ExitOnError)
	scope := fs.String("scope", "project", "Installation scope")
	_ = fs.Parse(args)

	skills, err := skill.ListInstalledSkills(*scope)
	if err != nil {
		return err
	}

	if len(skills) == 0 {
		fmt.Println("No skills installed.")
		return nil
	}

	fmt.Printf("%-20s %-15s %-30s %-20s\n", "NAME", "AGENT", "SOURCE", "INSTALLED")
	for _, s := range skills {
		fmt.Printf("%-20s %-15s %-30s %-20s\n", s.Meta.Name, s.Agent, s.Meta.OriginalSource, s.Meta.InstallTime.Format("2006-01-02"))
	}
	return nil
}

// RunSkillInspect is a subcommand `rntocase skill inspect` -- Inspect a skill
func RunSkillInspect(args []string) error {
	fs := flag.NewFlagSet("skill inspect", flag.ExitOnError)
	scope := fs.String("scope", "project", "Installation scope")
	agent := fs.String("agent", "common", "Target agent")
	outputJSON := fs.Bool("json", false, "Output in JSON format")
	_ = fs.Parse(args)

	positionalArgs := fs.Args()
	if len(positionalArgs) < 1 {
		return fmt.Errorf("usage: skill inspect <name>")
	}

	name := positionalArgs[0]

	meta, destDir, err := skill.InspectSkill(name, *scope, *agent)
	if err != nil {
		return err
	}

	if *outputJSON {
		data, _ := json.MarshalIndent(meta, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("Name:       %s\n", meta.Name)
	fmt.Printf("Path:       %s\n", destDir)
	fmt.Printf("Source:     %s\n", meta.OriginalSource)
	if meta.OwnerRepo != "" {
		fmt.Printf("Upstream:   %s\n", meta.OwnerRepo)
		fmt.Printf("Revision:   %s\n", meta.SourceRevision)
	}
	fmt.Printf("Installed:  %s\n", meta.InstallTime.Format(time.RFC3339))

	return nil
}
