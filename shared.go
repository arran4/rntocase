package rntocase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"path/filepath"
	"strings"
)

var inputReader *bufio.Reader

func getReader() *bufio.Reader {
	if inputReader == nil {
		inputReader = bufio.NewReader(os.Stdin)
	}
	return inputReader
}

// Confirm prompts the user with a yes/no question.
func Confirm(prompt string) bool {
	return ConfirmWithReader(prompt, getReader())
}

// SplitExtension splits a filename into its base name and extension.
// It supports compound extensions like .tar.gz and handles dotfiles cleanly.
func SplitExtension(filename string) (name, ext string) {
	if filename == "" || filename == "." || filename == ".." {
		return filename, ""
	}

	compoundExts := []string{".tar.gz", ".tar.bz2", ".tar.xz"}
	lower := strings.ToLower(filename)
	for _, ce := range compoundExts {
		if strings.HasSuffix(lower, ce) {
			name = filename[:len(filename)-len(ce)]
			ext = filename[len(filename)-len(ce):]
			if name == "" {
				return filename, ""
			}
			return name, ext
		}
	}

	ext = filepath.Ext(filename)
	name = filename[:len(filename)-len(ext)]

	if name == "" && ext != "" {
		return ext, ""
	}

	return name, ext
}

// RenameStatus represents the status of a rename operation
type RenameStatus string

const (
	StatusPlanned   RenameStatus = "planned"
	StatusRenamed   RenameStatus = "renamed"
	StatusUnchanged RenameStatus = "unchanged"
	StatusSkipped   RenameStatus = "skipped"
	StatusCollision RenameStatus = "collision"
	StatusFailed    RenameStatus = "failed"
)

// RenameOperation represents a single file rename operation in a structured format
type RenameOperation struct {
	Source      string       `json:"source"`
	Destination string       `json:"destination"`
	Status      RenameStatus `json:"status"`
	Error       string       `json:"error,omitempty"`
}

// RenameSummary contains aggregate counts for a rename batch
type RenameSummary struct {
	Planned   int `json:"planned"`
	Renamed   int `json:"renamed"`
	Unchanged int `json:"unchanged"`
	Skipped   int `json:"skipped"`
	Collision int `json:"collision"`
	Failed    int `json:"failed"`
}

// RenameResult represents the comprehensive result of a rename operation batch
type RenameResult struct {
	DryRun     bool              `json:"dry_run"`
	Operations []RenameOperation `json:"operations"`
	Summary    RenameSummary     `json:"summary"`
}

// RenamePlan represents a planned rename operation
type RenamePlan struct {
	OriginalPath string
	NewPath      string
	WillChange   bool
	Status       RenameStatus
	Error        error
}

// RenameError represents an error that occurred during the rename process
type RenameError struct {
	Path    string
	NewPath string
	Err     error
}

func (e *RenameError) Unwrap() error {
	return e.Err
}

func (e *RenameError) Error() string {
	return fmt.Sprintf("error processing '%s': %v", e.Path, e.Err)
}

// RenameBatchError is an aggregation of multiple RenameErrors
type RenameBatchError struct {
	Errors []*RenameError
}

func (e *RenameBatchError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("%d rename operations failed:\n  %s", len(e.Errors), strings.Join(msgs, "\n  "))
}

// ConfirmWithReader prompts the user with a yes/no question using the provided reader.
func ConfirmWithReader(prompt string, reader *bufio.Reader) bool {
	for {
		fmt.Print(prompt + " [y/n]: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Assuming 'n'.")
			return false
		}
		response = strings.TrimSpace(strings.ToLower(response))
		switch response {
		case "y":
			return true
		case "n":
			return false
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
}

// RenameFiles applies a renaming function to a list of files.
// Supports dry-run and interactive modes.
func RenameFiles(files []string, renameFunc func(string) (string, error), dryRun bool, interactive bool, outputJSON bool) error {
	if interactive && outputJSON {
		return fmt.Errorf("cannot use interactive mode with JSON output")
	}

	result := RenameResult{
		DryRun:     dryRun,
		Operations: []RenameOperation{},
		Summary:    RenameSummary{},
	}

	var plans []RenamePlan
	destCount := make(map[string]int)

	// Phase 1: Planning (Preflight)
	for _, file := range files {
		dir := filepath.Dir(file)
		base := filepath.Base(file)

		name, ext := SplitExtension(base)
		baseNameRenamed, err := renameFunc(name)
		if err != nil {
			plans = append(plans, RenamePlan{
				OriginalPath: file,
				Status:       StatusFailed,
				Error:        fmt.Errorf("failed to generate new name: %w", err),
			})
			continue
		}

		newName := baseNameRenamed + ext
		newPath := filepath.Join(dir, newName)

		absDir, err := filepath.Abs(dir)
		var evalDir string
		if err == nil {
			evalDir, _ = filepath.EvalSymlinks(absDir)
			if evalDir == "" {
				evalDir = absDir
			}
		} else {
			evalDir = dir
		}

		newCanonicalPath := filepath.Join(evalDir, newName)
		destKey := strings.ToLower(newCanonicalPath)

		plan := RenamePlan{
			OriginalPath: file,
			NewPath:      newPath,
			WillChange:   newPath != file,
			Status:       StatusPlanned,
		}

		if !plan.WillChange {
			plan.Status = StatusUnchanged
		} else {
			destCount[destKey]++
			if destCount[destKey] > 1 {
				plan.Error = fmt.Errorf("collision: multiple source files map to destination '%s'", newPath)
				plan.Status = StatusCollision

				// Retroactively flag the FIRST mapped file as a collision using the exact canonical key structure
				for i := range plans {
					// re-derive the peer's canonical key
					peerDir := filepath.Dir(plans[i].NewPath)
					peerBase := filepath.Base(plans[i].NewPath)
					peerAbsDir, err := filepath.Abs(peerDir)
					var peerEvalDir string
					if err == nil {
						peerEvalDir, _ = filepath.EvalSymlinks(peerAbsDir)
						if peerEvalDir == "" {
							peerEvalDir = peerAbsDir
						}
					} else {
						peerEvalDir = peerDir
					}
					peerCanonicalPath := filepath.Join(peerEvalDir, peerBase)
					peerDestKey := strings.ToLower(peerCanonicalPath)

					if peerDestKey == destKey && plans[i].Error == nil {
						plans[i].Error = fmt.Errorf("collision: multiple source files map to destination '%s'", plans[i].NewPath)
						plans[i].Status = StatusCollision
					}
				}
			} else {
				destStat, destErr := os.Lstat(newPath)
				if destErr == nil {
					srcStat, srcErr := os.Lstat(file)

					if srcErr != nil || !os.SameFile(srcStat, destStat) {
						plan.Error = fmt.Errorf("collision: destination '%s' already exists", newPath)
						plan.Status = StatusCollision
					}
				} else if !os.IsNotExist(destErr) {
					plan.Error = fmt.Errorf("collision: could not read destination '%s': %v", newPath, destErr)
					plan.Status = StatusCollision
				}
			}
		}

		plans = append(plans, plan)
	}

	var batchErrors []*RenameError

	// Ensure retroactive errors are collected in batchErrors
	// Phase 1b: Batch Safety Validation
	for _, plan := range plans {
		if plan.Error != nil {
			if !outputJSON {
				fmt.Fprintf(os.Stderr, "Error planning '%s': %v\n", plan.OriginalPath, plan.Error)
			}
			batchErrors = append(batchErrors, &RenameError{
				Path:    plan.OriginalPath,
				NewPath: plan.NewPath,
				Err:     plan.Error,
			})
		}
	}

	if len(batchErrors) > 0 {
		// If batch fails, write ALL plans to the result before aborting
		for _, plan := range plans {
			if plan.Error != nil {
				if plan.Status == StatusCollision {
					result.Summary.Collision++
				} else {
					plan.Status = StatusFailed
					result.Summary.Failed++
				}
				result.Operations = append(result.Operations, RenameOperation{
					Source:      plan.OriginalPath,
					Destination: plan.NewPath,
					Status:      plan.Status,
					Error:       plan.Error.Error(),
				})
			} else {
				// Record successful plans as skipped (or unchanged if they wouldn't have changed)
				status := StatusSkipped
				if plan.Status == StatusUnchanged {
					status = StatusUnchanged
					result.Summary.Unchanged++
				} else {
					result.Summary.Skipped++
				}
				result.Operations = append(result.Operations, RenameOperation{
					Source:      plan.OriginalPath,
					Destination: plan.NewPath,
					Status:      status,
					Error:       "batch aborted due to other errors",
				})
			}
		}

		if outputJSON {
			jsonBytes, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(jsonBytes))
		}
		return &RenameBatchError{Errors: batchErrors}
	}

	// Phase 2: Execution & Reporting
	for _, plan := range plans {
		if !plan.WillChange {
			if !outputJSON {
				fmt.Printf("Skipping '%s' (already matches desired format).\n", plan.OriginalPath)
			}
			result.Summary.Unchanged++
			result.Operations = append(result.Operations, RenameOperation{
				Source:      plan.OriginalPath,
				Destination: plan.NewPath,
				Status:      StatusUnchanged,
			})
			continue
		}

		if dryRun {
			if !outputJSON {
				fmt.Printf("Rename: '%s' -> '%s'\n", plan.OriginalPath, plan.NewPath)
			}
			result.Summary.Planned++
			result.Operations = append(result.Operations, RenameOperation{
				Source:      plan.OriginalPath,
				Destination: plan.NewPath,
				Status:      StatusPlanned,
			})
			continue
		}

		if interactive {
			if !Confirm(fmt.Sprintf("Rename '%s' to '%s'?", plan.OriginalPath, filepath.Base(plan.NewPath))) {
				if !outputJSON {
					fmt.Println("Skipped.")
				}
				result.Summary.Skipped++
				result.Operations = append(result.Operations, RenameOperation{
					Source:      plan.OriginalPath,
					Destination: plan.NewPath,
					Status:      StatusSkipped,
				})
				continue
			}
		}

		if !outputJSON {
			fmt.Printf("Rename: '%s' -> '%s'\n", plan.OriginalPath, plan.NewPath)
		}

		if err := os.Rename(plan.OriginalPath, plan.NewPath); err != nil {
			if !outputJSON {
				fmt.Fprintf(os.Stderr, "Error renaming '%s': %v\n", plan.OriginalPath, err)
			}
			batchErrors = append(batchErrors, &RenameError{
				Path:    plan.OriginalPath,
				NewPath: plan.NewPath,
				Err:     err,
			})
			result.Summary.Failed++
			result.Operations = append(result.Operations, RenameOperation{
				Source:      plan.OriginalPath,
				Destination: plan.NewPath,
				Status:      StatusFailed,
				Error:       err.Error(),
			})
			continue
		}

		if !outputJSON {
			fmt.Println("Renamed successfully.")
		}
		result.Summary.Renamed++
		result.Operations = append(result.Operations, RenameOperation{
			Source:      plan.OriginalPath,
			Destination: plan.NewPath,
			Status:      StatusRenamed,
		})
	}

	if outputJSON {
		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonBytes))
	}

	if len(batchErrors) > 0 {
		return &RenameBatchError{Errors: batchErrors}
	}

	return nil
}

// LoadAcronymsFromFile reads acronyms from a file and configures them.
func LoadAcronymsFromFile(filePath string) error {
	return fmt.Errorf("acronym configuration is no longer globally supported")
}

func Run(converter func(string) (string, error), value string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	result, err = converter(value)
	return result, err
}

func RenderUsageTable(converter func(string) (string, error)) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Input", "Output"})
	for _, values := range ExampleGroups {
		tw.AppendRow(table.Row{"", values.Name})
		for _, value := range values.Values {
			result, err := Run(converter, value)
			if err != nil {
				result = "!!!Error!!!"
			}
			tw.AppendRow(table.Row{value, result})
		}
	}
	tw.SetOutputMirror(os.Stderr)
	tw.Render()
}
