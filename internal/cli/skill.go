package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arran4/rntocase/internal/skill"
)

// RunSkill is a subcommand `rntocase skill` -- Manage AI agent skills for this CLI
func RunSkill() error {
	return fmt.Errorf("skill command requires a subcommand (install, update, remove, list, inspect)")
}

// RunSkillInstall is a subcommand `rntocase skill install` -- Install a skill
//
// Flags:
//
//	scope:      --scope (default: "project") Installation scope: user or project
//	agent:      --agent (default: "common") Target agent
//	replace:    --replace Replace existing skill
//	ref:        --ref Specific Git ref
//	path:       --path Repository subdirectory
//	name:       --name Local skill name
//	source:     @1 (min: 1) Source repository/path
//	nameOrPath: @2 (default: "") Optional fallback name/path
//
// Examples:
//
//	rntocase skill install --ref v1.2.0 --path skills/example --name example owner/repo
func RunSkillInstall(
	scope string,
	agent string,
	replace bool,
	ref string,
	path string,
	name string,
	source string,
	nameOrPath string,
) error {
	if scope == "" {
		scope = "project"
	}
	if agent == "" {
		agent = "common"
	}

	target, err := skill.ResolveTarget(scope, agent)
	if err != nil {
		return err
	}

	isLocal, isOfficial, ownerRepo, err := skill.ClassifySource(source)
	if err != nil {
		return err
	}

	// Mapping old positional arguments to new concepts
	pathWithin := path
	skillName := name
	if nameOrPath != "" {
		if skillName == "" {
			if !isLocal && !isOfficial && strings.Contains(nameOrPath, "/") {
				if pathWithin == "" {
					pathWithin = nameOrPath
				}
			} else {
				skillName = nameOrPath
			}
		} else if pathWithin == "" && !isLocal && !isOfficial && strings.Contains(nameOrPath, "/") {
			pathWithin = nameOrPath
		}
	}

	// Determine final skill name
	if skillName == "" {
		if isLocal {
			skillName = filepath.Base(source)
		} else if isOfficial {
			skillName = "rntocase"
		} else {
			if pathWithin != "" {
				skillName = filepath.Base(pathWithin)
			} else {
				parts := strings.Split(ownerRepo, "/")
				if len(parts) >= 2 {
					skillName = parts[1] // fallback to repo name
				} else {
					return fmt.Errorf("could not determine skill name automatically, please provide it using --name")
				}
			}
		}
	}

	// Validate skillName doesn't escape the target dir or overwrite it directly
	destDir, err := skill.ResolveSkillPath(target, skillName)
	if err != nil {
		return fmt.Errorf("invalid skill name: %w", err)
	}

	// Check if destination exists before proceeding
	if _, err := os.Stat(destDir); err == nil {
		if !replace {
			if _, metaErr := skill.LoadMetadata(destDir); metaErr == nil {
				return fmt.Errorf("skill '%s' is already installed at %s. Use 'skill update' to update it, or use --replace to force reinstall", skillName, destDir)
			}
			return fmt.Errorf("destination %s already exists. Use --replace to overwrite", destDir)
		}
	}

	fmt.Printf("Installing skill '%s' to %s...\n", skillName, destDir)

	meta := &skill.Metadata{
		Name:           skillName,
		OriginalSource: source,
		InstallTime:    time.Now(),
		InstallerApp:   "rntocase",
	}

	var tarPath string

	if !isLocal && !isOfficial {
		// Remote Github Download outside ReplaceSafely to minimize staging time
		var sha string
		var err error
		tarPath, sha, err = skill.DownloadGitHubRepository(ownerRepo, ref)
		if err != nil {
			return fmt.Errorf("failed to download skill: %w", err)
		}
		defer func() { _ = os.Remove(tarPath) }()

		meta.OwnerRepo = ownerRepo
		meta.SourceRevision = sha
		meta.PathWithin = pathWithin
		meta.RequestedRef = ref
	}

	err = skill.ReplaceSafely(destDir, func(stagingDir string) error {
		if isLocal {
			if err := skill.CopyLocalDirectory(source, stagingDir); err != nil {
				return fmt.Errorf("failed to copy local skill: %w", err)
			}
		} else if isOfficial {
			if err := skill.ExtractEmbeddedSkill("rntocase", stagingDir); err != nil {
				return fmt.Errorf("failed to install official skill: %w", err)
			}
		} else {
			if err := skill.ExtractTarGz(tarPath, stagingDir, pathWithin); err != nil {
				return fmt.Errorf("failed to extract skill: %w", err)
			}
		}

		// Validate SKILL.md
		skillMdPath := filepath.Join(stagingDir, "SKILL.md")
		mdContent, err := os.ReadFile(skillMdPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("installation failed: skill must contain a SKILL.md file")
			}
			return fmt.Errorf("failed to read SKILL.md: %w", err)
		}

		manifest, err := skill.ParseAndValidateManifest(mdContent)
		if err != nil {
			return fmt.Errorf("invalid SKILL.md manifest: %w", err)
		}

		if manifest.Name != meta.Name {
			return fmt.Errorf("installation failed: skill name in manifest ('%s') does not match installed directory name ('%s')", manifest.Name, meta.Name)
		}

		digest, err := skill.ComputeDirectoryDigest(stagingDir)
		if err == nil {
			meta.ContentDigest = digest
		}

		if err := skill.SaveMetadata(stagingDir, meta); err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("Skill installed successfully.")
	return nil
}

// RunSkillUpdate is a subcommand `rntocase skill update` -- Update a skill
// Flags:
//
//	scope: --scope (default: "project") Installation scope
//	agent: --agent (default: "common") Target agent
//	force: --force Force update and overwrite local changes
//	all:   --all Update all installed skills in the given scope
//	name:  @1 The name of the skill to update
func RunSkillUpdate(scope string, agent string, force bool, all bool, name string) error {
	if scope == "" {
		scope = "project"
	}
	if agent == "" {
		agent = "common"
	}

	if name == "" && !all {
		return fmt.Errorf("usage: skill update <name> or skill update --all")
	}

	type skillUpdateTask struct {
		Name string
		Meta *skill.Metadata
		Dir  string
		Err  error
	}

	var skillsToUpdate []skillUpdateTask

	if all {
		installed, err := skill.ListInstalledSkills(scope, agent)
		if err != nil {
			return err
		}
		for _, info := range installed {
			// Find the actual path using the agent it was found under
			meta, dir, err := skill.InspectSkill(info.Meta.Name, scope, info.Agent)
			if err == nil {
				skillsToUpdate = append(skillsToUpdate, skillUpdateTask{
					Name: info.Meta.Name,
					Meta: meta,
					Dir:  dir,
				})
			} else {
				skillsToUpdate = append(skillsToUpdate, skillUpdateTask{
					Name: info.Meta.Name,
					Err:  err,
				})
			}
		}
	} else {
		meta, destDir, err := skill.InspectSkill(name, scope, agent)
		if err != nil {
			// for single skill update, just return the inspect error immediately
			return err
		}

		// explicitly preserve pre-existing single-skill behavior for local-only skills
		if meta.OwnerRepo == "" && meta.OriginalSource != "official" && meta.OriginalSource != "rntocase" {
			return fmt.Errorf("skill '%s' is locally installed and cannot be updated automatically", name)
		}

		skillsToUpdate = append(skillsToUpdate, skillUpdateTask{
			Name: name,
			Meta: meta,
			Dir:  destDir,
		})
	}

	if len(skillsToUpdate) == 0 {
		fmt.Println("No skills found to update.")
		return nil
	}

	type updateResult struct {
		Name   string
		Status string
		Err    error
	}
	var results []updateResult

	for _, s := range skillsToUpdate {
		if s.Err != nil {
			results = append(results, updateResult{Name: s.Name, Status: "inspection failed", Err: s.Err})
			continue
		}

		status, err := updateSingleSkill(s.Name, s.Meta, s.Dir, force)
		results = append(results, updateResult{Name: s.Name, Status: status, Err: err})
	}

	fmt.Println("\nUpdate Summary:")
	fmt.Printf("%-20s %-40s %s\n", "NAME", "STATUS", "ERROR")
	fmt.Println(strings.Repeat("-", 80))

	var failedSkills []string
	for _, res := range results {
		errStr := ""
		if res.Err != nil {
			errStr = res.Err.Error()
			failedSkills = append(failedSkills, fmt.Sprintf("%s (%s: %v)", res.Name, res.Status, res.Err))
		} else if res.Status == "update failed" || res.Status == "inspection failed" {
			failedSkills = append(failedSkills, fmt.Sprintf("%s (%s)", res.Name, res.Status))
		}

		fmt.Printf("%-20s %-40s %s\n", res.Name, res.Status, errStr)
	}
	fmt.Println()

	if len(failedSkills) > 0 {
		return fmt.Errorf("one or more skills failed to update:\n- %s", strings.Join(failedSkills, "\n- "))
	}

	return nil
}

func updateSingleSkill(name string, meta *skill.Metadata, destDir string, force bool) (string, error) {
	return updateSingleSkillWithExtractor(name, meta, destDir, force, skill.ExtractEmbeddedSkill)
}

func updateSingleSkillWithExtractor(name string, meta *skill.Metadata, destDir string, force bool, extractEmbedded func(assetName, destDir string) error) (string, error) {
	if meta.OriginalSource == "official" || meta.OriginalSource == "rntocase" {
		fmt.Printf("Checking for updates for official skill '%s'...\n", name)

		// If there is an update to the binary it could bring a new embedded skill
		// Since we don't have a network version for embedded, we just forcefully re-extract it
		// Check for local mods first
		currentDigest, err := skill.ComputeDirectoryDigest(destDir)
		if err == nil && meta.ContentDigest != "" && currentDigest != meta.ContentDigest && !force {
			fmt.Printf("Skill '%s' has local modifications. Use --force to replace.\n", name)
			return "skipped because of local modifications", nil
		}

		err = skill.ReplaceSafely(destDir, func(stagingDir string) error {
			if err := extractEmbedded("rntocase", stagingDir); err != nil {
				return err
			}

			// Validate SKILL.md for official/embedded update too
			skillMdPath := filepath.Join(stagingDir, "SKILL.md")
			mdContent, err := os.ReadFile(skillMdPath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("update failed: new skill version must contain a SKILL.md file")
				}
				return fmt.Errorf("failed to read SKILL.md: %w", err)
			}

			manifest, err := skill.ParseAndValidateManifest(mdContent)
			if err != nil {
				return fmt.Errorf("invalid SKILL.md manifest: %w", err)
			}

			if manifest.Name != meta.Name {
				return fmt.Errorf("update failed: skill name in manifest ('%s') does not match installed directory name ('%s')", manifest.Name, meta.Name)
			}

			digest, _ := skill.ComputeDirectoryDigest(stagingDir)
			meta.ContentDigest = digest
			meta.InstallTime = time.Now()
			if err := skill.SaveMetadata(stagingDir, meta); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return "update failed", err
		}

		fmt.Printf("Skill '%s' updated from embedded source.\n", name)
		return "updated", nil
	}

	if meta.OwnerRepo == "" {
		fmt.Printf("Skill '%s' is locally installed and cannot be updated automatically\n", name)
		return "unsupported/local-only", nil
	}

	fmt.Printf("Checking for updates for '%s'...\n", name)
	hasUpdate, shaUpdate, err := skill.CheckUpdate(meta)
	if err != nil {
		return "update failed", fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Println("Skill is already up to date.")
		return "already current", nil
	}

	// Check for local modifications
	currentDigest, err := skill.ComputeDirectoryDigest(destDir)
	if err == nil && meta.ContentDigest != "" && currentDigest != meta.ContentDigest && !force {
		fmt.Printf("Skill '%s' has local modifications. Use --force to replace.\n", name)
		return "skipped because of local modifications", nil
	}

	fmt.Println("Updating skill...")

	// Re-download first to minimize staging time
	tarPath, shaDownload, err := skill.DownloadGitHubRepository(meta.OwnerRepo, meta.RequestedRef)
	if err != nil {
		return "update failed", fmt.Errorf("failed to download update: %w", err)
	}
	defer func() { _ = os.Remove(tarPath) }()

	err = skill.ReplaceSafely(destDir, func(stagingDir string) error {
		if err := skill.ExtractTarGz(tarPath, stagingDir, meta.PathWithin); err != nil {
			return fmt.Errorf("failed to extract updated skill: %w", err)
		}

		// Validate SKILL.md for update too
		skillMdPath := filepath.Join(stagingDir, "SKILL.md")
		mdContent, err := os.ReadFile(skillMdPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("update failed: new skill version must contain a SKILL.md file")
			}
			return fmt.Errorf("failed to read SKILL.md: %w", err)
		}

		manifest, err := skill.ParseAndValidateManifest(mdContent)
		if err != nil {
			return fmt.Errorf("invalid SKILL.md manifest: %w", err)
		}

		if manifest.Name != meta.Name {
			return fmt.Errorf("update failed: skill name in manifest ('%s') does not match installed directory name ('%s')", manifest.Name, meta.Name)
		}

		meta.SourceRevision = shaDownload
		meta.InstallTime = time.Now()

		digest, err := skill.ComputeDirectoryDigest(stagingDir)
		if err == nil {
			meta.ContentDigest = digest
		}

		if err := skill.SaveMetadata(stagingDir, meta); err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}

		return nil
	})

	if err != nil {
		return "update failed", err
	}

	shortSha := shaUpdate
	if len(shortSha) > 7 {
		shortSha = shortSha[:7]
	}

	fmt.Printf("Skill updated successfully to revision %s.\n", shortSha)
	return "updated", nil
}

// RunSkillRemove is a subcommand `rntocase skill remove` -- Remove a skill
//
// Flags:
//
//	scope: --scope (default: "project") Installation scope
//	agent: --agent (default: "common") Target agent
//	name:  @1 (min: 1) The name of the skill to remove
func RunSkillRemove(scope string, agent string, name string) error {
	if scope == "" {
		scope = "project"
	}
	if agent == "" {
		agent = "common"
	}

	if name == "" {
		return fmt.Errorf("usage: skill remove <name>")
	}

	if err := skill.RemoveSkill(name, scope, agent); err != nil {
		return err
	}

	fmt.Printf("Skill '%s' removed successfully.\n", name)
	return nil
}

// RunSkillList is a subcommand `rntocase skill list` -- List skills
//
// Flags:
//
//	scope: --scope (default: "project") Installation scope
func RunSkillList(scope string) error {
	if scope == "" {
		scope = "project"
	}

	skills, err := skill.ListInstalledSkills(scope, "all")
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
//
// Flags:
//
//	scope:      --scope (default: "project") Installation scope
//	agent:      --agent (default: "common") Target agent
//	outputJSON: --json Output in JSON format
//	name:       @1 (min: 1) The name of the skill to inspect
func RunSkillInspect(scope string, agent string, outputJSON bool, name string) error {
	if scope == "" {
		scope = "project"
	}
	if agent == "" {
		agent = "common"
	}

	if name == "" {
		return fmt.Errorf("usage: skill inspect <name>")
	}

	meta, destDir, err := skill.InspectSkill(name, scope, agent)
	if err != nil {
		return err
	}

	if outputJSON {
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
