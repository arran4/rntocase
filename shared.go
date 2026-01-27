package rntocase

import (
	"bufio"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jedib0t/go-pretty/table"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// Confirm prompts the user with a yes/no question.
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt + " [y/n]: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Assuming 'n'.")
			return false
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" {
			return true
		} else if response == "n" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
}

// RenameFiles applies a renaming function to a list of files.
// Supports dry-run and interactive modes.
func RenameFiles(files []string, renameFunc func(string) (string, error), dryRun bool, interactive bool) error {
	for _, file := range files {
		// Extract path and filename
		dir := filepath.Dir(file)
		base := filepath.Base(file)

		// Generate the new filename
		baseNameRenamed, err := renameFunc(strings.TrimSuffix(base, filepath.Ext(base)))
		if err != nil {
			return err
		}
		newName := baseNameRenamed + filepath.Ext(base)
		newPath := filepath.Join(dir, newName)

		// Skip if no changes
		if newPath == file {
			fmt.Printf("Skipping '%s' (already matches desired format).\n", file)
			continue
		}

		// Display the intended change
		fmt.Printf("Rename: '%s' -> '%s'\n", file, newPath)

		// Dry-run mode
		if dryRun {
			continue
		}

		// Interactive mode
		if interactive {
			if !Confirm(fmt.Sprintf("Rename '%s' to '%s'?", file, newName)) {
				fmt.Println("Skipped.")
				continue
			}
		}

		// Perform the rename
		if err := os.Rename(file, newPath); err != nil {
			fmt.Printf("Error renaming '%s': %v\n", file, err)
			continue
		}
		fmt.Println("Renamed successfully.")
	}
	return nil
}

// LoadAcronymsFromFile reads acronyms from a file and configures them for strcase.
func LoadAcronymsFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open acronym file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if word == "" {
			continue // Skip empty lines
		}
		strcase.ConfigureAcronym(word, word) // Configure acronym to retain its casing
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading acronym file: %w", err)
	}

	fmt.Printf("Loaded acronyms from file: %s\n", filePath)
	return nil
}

func Run(algos map[string]func(string) (string, error), key string, value string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	result, err = algos[key](value)
	return result, err
}

func RenderUsageTable(algos map[string]func(string) (string, error)) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Algorithm"})
	keys := slices.Collect(maps.Keys(algos))
	sort.Strings(keys)
	for _, values := range ExampleGroups {
		tw.AppendRow(slices.Collect(func(yield func(row any) bool) {
			yield("")
			if !yield(values.Name) {
				return
			}
			for _, value := range values.Values {
				if !yield(value) {
					return
				}
			}
		}))
		for _, key := range keys {
			tw.AppendRow(slices.Collect(func(yield func(row any) bool) {
				yield(key)
				yield("")
				runSafeBatch(algos[key], values.Values, func(s string) bool {
					return yield(s)
				})
			}))
		}
	}
	tw.SetOutputMirror(os.Stderr)
	tw.Render()
}

// runSafeBatch executes the algo function for each value in the slice safely.
// It ensures that panic in one item does not stop processing of subsequent items,
// but minimizes defer overhead by using a single defer per batch (recursively if needed).
func runSafeBatch(algo func(string) (string, error), values []string, yield func(string) bool) {
	processBatch(algo, values, 0, yield)
}

func processBatch(algo func(string) (string, error), values []string, start int, yield func(string) bool) {
	var current int
	defer func() {
		if r := recover(); r != nil {
			// If panic occurs, we yield error for the current item
			if !yield("!!!Error!!!") {
				return
			}
			// Resume processing subsequent items
			if current+1 < len(values) {
				processBatch(algo, values, current+1, yield)
			}
		}
	}()

	for i := start; i < len(values); i++ {
		current = i
		res, err := algo(values[i])
		var val string
		if err != nil {
			val = "!!!Error!!!"
		} else {
			val = res
		}
		if !yield(val) {
			return
		}
	}
}
