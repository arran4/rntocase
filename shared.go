package rntocase

import (
	"bufio"
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

// RenamePlan represents a planned rename operation
type RenamePlan struct {
	OriginalPath string
	NewPath      string
	WillChange   bool
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
func RenameFiles(files []string, renameFunc func(string) (string, error), dryRun bool, interactive bool) error {
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
				Error:        fmt.Errorf("failed to generate new name: %w", err),
			})
			continue
		}

		newName := baseNameRenamed + ext
		newPath := filepath.Join(dir, newName)

		plan := RenamePlan{
			OriginalPath: file,
			NewPath:      newPath,
			WillChange:   newPath != file,
		}

		if plan.WillChange {
			// Check if duplicate destination within batch
			destCount[newPath]++
			if destCount[newPath] > 1 {
				plan.Error = fmt.Errorf("collision: multiple source files map to destination '%s'", newPath)
			} else {
				// Check if destination exists (and isn't just a case-change rename of the same file)
				// A case-change on a case-insensitive FS will mean os.Stat returns no error,
				// but os.SameFile won't work on non-existent files.
				// We can just check os.Stat and if it exists and is not the same file
				if destStat, err := os.Stat(newPath); err == nil {
					srcStat, srcErr := os.Stat(file)
					if srcErr != nil || !os.SameFile(srcStat, destStat) {
						plan.Error = fmt.Errorf("collision: destination '%s' already exists", newPath)
					}
				}
			}
		}

		plans = append(plans, plan)
	}

	var batchErrors []*RenameError

	// Phase 1b: Batch Safety Validation
	// If any planning error occurred (e.g. collision), abort the entire batch
	// to ensure we don't leave the filesystem in a partially mutated state.
	for _, plan := range plans {
		if plan.Error != nil {
			fmt.Printf("Error planning '%s': %v\n", plan.OriginalPath, plan.Error)
			batchErrors = append(batchErrors, &RenameError{
				Path:    plan.OriginalPath,
				NewPath: plan.NewPath,
				Err:     plan.Error,
			})
		}
	}

	if len(batchErrors) > 0 {
		return &RenameBatchError{Errors: batchErrors}
	}

	// Phase 2: Execution & Reporting
	for _, plan := range plans {
		if !plan.WillChange {
			fmt.Printf("Skipping '%s' (already matches desired format).\n", plan.OriginalPath)
			continue
		}

		fmt.Printf("Rename: '%s' -> '%s'\n", plan.OriginalPath, plan.NewPath)

		if dryRun {
			continue
		}

		if interactive {
			if !Confirm(fmt.Sprintf("Rename '%s' to '%s'?", plan.OriginalPath, filepath.Base(plan.NewPath))) {
				fmt.Println("Skipped.")
				continue
			}
		}

		if err := os.Rename(plan.OriginalPath, plan.NewPath); err != nil {
			fmt.Printf("Error renaming '%s': %v\n", plan.OriginalPath, err)
			batchErrors = append(batchErrors, &RenameError{
				Path:    plan.OriginalPath,
				NewPath: plan.NewPath,
				Err:     err,
			})
			continue
		}
		fmt.Println("Renamed successfully.")
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
