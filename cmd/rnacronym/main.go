package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/arran4/strings2"
	"os"
	"strings"
)

const (
	appName     = "rnacronym"
)

func converter(s string) (string, error) {
    // Rely completely on string2.Parse output to extract acronym chars.
	words, err := strings2.Parse(s)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for _, w := range words {
		for _, r := range w.String() {
			result.WriteString(strings.ToUpper(string(r)))
			break
		}
	}
	return result.String(), nil
}

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: "+appName+" [options] <file1> [<file2> ...]")
		_, _ = fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Conversion Examples:")
		rntocase.RenderUsageTable(converter)
	}

	flag.Parse()

	// Get file arguments
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Error: No files provided.")
		flag.Usage()
		os.Exit(1)
	}

	// Use the shared RenameFiles function
	if err := rntocase.RenameFiles(files, converter, *dryRun, *interactive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
