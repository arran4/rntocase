package skill

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed embedded_skills/*
var EmbeddedSkillsFS embed.FS

// ExtractEmbeddedSkill extracts a skill from the embedded filesystem to the destination directory.
func ExtractEmbeddedSkill(skillName, destDir string) error {
	basePath := "embedded_skills/" + skillName

	err := fs.WalkDir(EmbeddedSkillsFS, basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(destDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := fs.ReadFile(EmbeddedSkillsFS, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, content, 0644)
	})

	if err != nil {
		if os.IsNotExist(err) || err == fs.ErrNotExist {
			return fmt.Errorf("embedded skill '%s' not found", skillName)
		}
		return err
	}

	return nil
}
