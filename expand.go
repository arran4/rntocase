package rntocase

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// ExpandFiles discovers and filters files based on recursion and include/exclude patterns.
func ExpandFiles(files []string, recursive bool, includes []string, excludes []string) ([]string, error) {
	if !recursive {
		return files, nil
	}

	// Validate patterns early
	for _, p := range includes {
		if !doublestar.ValidatePattern(p) {
			return nil, fmt.Errorf("invalid include pattern: %s", p)
		}
	}
	for _, p := range excludes {
		if !doublestar.ValidatePattern(p) {
			return nil, fmt.Errorf("invalid exclude pattern: %s", p)
		}
	}

	discoveredMap := make(map[string]bool)
	var discovered []string

	for _, root := range files {
		rootStat, err := os.Lstat(root)
		if err != nil {
			if recursive {
				return nil, fmt.Errorf("recursive root not accessible: %w", err)
			}
			// Let it pass through to the RenameFiles step which will fail with a good error
			if !discoveredMap[root] {
				discoveredMap[root] = true
				discovered = append(discovered, root)
			}
			continue
		}

		if !rootStat.IsDir() {
			if recursive {
				// Apply filtering to the file root
				slashRelPath := filepath.Base(root)

				excluded := false
				for _, pattern := range excludes {
					matchPattern := pattern
					if !strings.HasPrefix(pattern, "/") && !strings.HasPrefix(pattern, "**") {
						matchPattern = "**/" + pattern
					}
					matched, _ := doublestar.Match(matchPattern, slashRelPath)
					if matched {
						excluded = true
						break
					}
					if !strings.Contains(pattern, "/") {
						matchedBase, _ := doublestar.Match(pattern, slashRelPath)
						if matchedBase {
							excluded = true
							break
						}
					}
				}

				if excluded {
					continue
				}

				included := len(includes) == 0
				if !included {
					for _, pattern := range includes {
						matchPattern := pattern
						if !strings.HasPrefix(pattern, "/") && !strings.HasPrefix(pattern, "**") {
							matchPattern = "**/" + pattern
						}
						matched, _ := doublestar.Match(matchPattern, slashRelPath)
						if matched {
							included = true
							break
						}
						if !strings.Contains(pattern, "/") {
							matchedBase, _ := doublestar.Match(pattern, slashRelPath)
							if matchedBase {
								included = true
								break
							}
						}
					}
				}

				if included {
					if !discoveredMap[root] {
						discoveredMap[root] = true
						discovered = append(discovered, root)
					}
				}
				continue
			}

			// It's a file, keep it (non-recursive)
			if !discoveredMap[root] {
				discoveredMap[root] = true
				discovered = append(discovered, root)
			}
			continue
		}

		// Handle directory walking if recursive
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Exclude the root dir itself from being returned as a file
			if path == root {
				return nil
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				relPath = path
			}

			// Normalize to slashes for matching
			slashRelPath := filepath.ToSlash(relPath)

			// Check exclusions
			excluded := false
			for _, pattern := range excludes {
				// Normalize pattern to handle `.git/**` -> `**/.git/**` internally if it doesn't start with `/` or `**`
				matchPattern := pattern
				if !strings.HasPrefix(pattern, "/") && !strings.HasPrefix(pattern, "**") {
					// We prepend **/ so it matches anywhere in the tree
					matchPattern = "**/" + pattern
				}

				matched, _ := doublestar.Match(matchPattern, slashRelPath)
				if matched {
					excluded = true
					break
				}

				// Handle implicit pruning: if the pattern is `**/.git/**`, we want to skip `sub/.git` itself,
				// not just `sub/.git/config`. `doublestar.Match("**/.git/**", "sub/.git")` might be false,
				// so if pattern ends with `/**`, also check without the `/**` to skip the directory early.
				if strings.HasSuffix(matchPattern, "/**") {
					dirPattern := strings.TrimSuffix(matchPattern, "/**")
					matchedDir, _ := doublestar.Match(dirPattern, slashRelPath)
					if matchedDir {
						excluded = true
						break
					}
				}

				// Also try base name if the original pattern doesn't contain a slash
				if !strings.Contains(pattern, "/") {
					matchedBase, _ := doublestar.Match(pattern, filepath.Base(slashRelPath))
					if matchedBase {
						excluded = true
						break
					}
				}
			}

			if excluded {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if info.IsDir() {
				// If it's a symlink directory, don't follow it
				if (info.Mode() & os.ModeSymlink) != 0 {
					return filepath.SkipDir
				}
				return nil // Don't rename dirs by default
			}

			// Check inclusions
			included := len(includes) == 0
			if !included {
				for _, pattern := range includes {
					matchPattern := pattern
					if !strings.HasPrefix(pattern, "/") && !strings.HasPrefix(pattern, "**") {
						matchPattern = "**/" + pattern
					}

					matched, _ := doublestar.Match(matchPattern, slashRelPath)
					if matched {
						included = true
						break
					}
					// Also try base name if the pattern doesn't contain a slash
					if !strings.Contains(pattern, "/") {
						matchedBase, _ := doublestar.Match(pattern, filepath.Base(slashRelPath))
						if matchedBase {
							included = true
							break
						}
					}
				}
			}

			if included {
				// Don't follow symlinked files when discovered recursively
				if (info.Mode() & os.ModeSymlink) != 0 {
					return nil
				}

				if !discoveredMap[path] {
					discoveredMap[path] = true
					discovered = append(discovered, path)
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	// Deterministic sort
	sort.Strings(discovered)
	return discovered, nil
}

// RenameFilesWithDiscovery wraps RenameFiles by first expanding the files list
// recursively with include/exclude filters.
func RenameFilesWithDiscovery(files []string, recursive bool, includes []string, excludes []string, renameFunc func(string) (string, error), dryRun bool, interactive bool, outputJSON bool) error {
	expandedFiles, err := ExpandFiles(files, recursive, includes, excludes)
	if err != nil {
		if outputJSON {
			// Best effort to emit JSON even on flag validation errors
			// But for simplicity of matching the existing pattern, just return error
			return err
		}
		return err
	}

	return RenameFiles(expandedFiles, renameFunc, dryRun, interactive, outputJSON)
}
