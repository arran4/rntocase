package rntocase

import (
	"bufio"
	"fmt"
	"github.com/iancoleman/strcase"
	"os"
	"path/filepath"
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
