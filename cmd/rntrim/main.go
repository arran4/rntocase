package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"strings"
)

const (
	appName     = "rntrim"
)

var trimChars *string

func converter(s string) (string, error) {
	if *trimChars == "" {
		return strings.TrimSpace(s), nil
	} else {
		return strings.Trim(s, *trimChars), nil
	}
}

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	trimChars = flag.String("trim", "", "Characters to trim off end and start of name white space if not set.")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: "+appName+" [options] <file1> [<file2> ...]")
		_, _ = fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Conversion Examples:")
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"Algorithm"})

		for _, values := range rntocase.ExampleGroups {
            // Replicating original table rendering logic without algos map
            row := make(table.Row, 0)
			row = append(row, "")
			row = append(row, values.Name)
            for _, value := range values.Values {
                row = append(row, value)
            }
			tw.AppendRow(row)

            resRow := make(table.Row, 0)
            resRow = append(resRow, "go")
            resRow = append(resRow, "")
            for _, value := range values.Values {
                result, err := converter(value)
                if err != nil {
                    panic(err)
                }
                resRow = append(resRow, result)
            }
			tw.AppendRow(resRow)
		}
		tw.SetOutputMirror(os.Stderr)
		tw.Render()
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
