package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/arran4/strings2"
	"os"
)

const (
	appName     = "rntodelimited"
)

var delimiter *string

func converter(s string) (string, error) {
	return strings2.ToFormattedString(s, strings2.OptionDelimiter(*delimiter), strings2.OptionCaseMode(strings2.CMWhispering))
}

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	delimiter = flag.String("delimiter", "_", "The delimiter string to separate words")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: "+appName+" [options] <file1> [<file2> ...]")
		_, _ = fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Conversion Examples:")
		rntocase.RenderUsageTable(map[string]func(string) (string, error){"strings2": converter})
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
