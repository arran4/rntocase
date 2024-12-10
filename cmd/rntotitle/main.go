package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/gobeam/stringy"
	"github.com/revett/titlecase"
	skstrings "github.com/searKing/golang/go/strings"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"strings"
)

const (
	defaultAlgo = "revett"
	caseType    = "titlecase"
	appName     = "rntotitle"
)

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	var opArgs []titlecase.Option
	var (
		algos = map[string]func(string) (string, error){
			"revett": func(s string) (string, error) {
				return titlecase.Convert(s, opArgs...)
			},
			"searking": func(s string) (string, error) {
				return skstrings.TitleCase(s), nil
			},
			"go": func(s string) (string, error) {
				return strings.ToTitle(s), nil
			},
			"resenje": func(s string) (string, error) {
				return casbab.Title(s), nil
			},
			"gobeam": func(s string) (string, error) {
				return stringy.New(s).Title(), nil
			},
		}
	)

	algorithm := flag.String("algorithm", defaultAlgo, "Choose the "+caseType+" algorithm to use, supported: "+strings.Join(slices.Collect(maps.Keys(algos)), ",")+".")
	flag.Usage = func() {
		fmt.Println("Usage: " + appName + " [options] <file1> [<file2> ...]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Get file arguments
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Error: No files provided.")
		flag.Usage()
		os.Exit(1)
	}

	converter, ok := algos[*algorithm]
	if !ok {
		fmt.Printf("Uunsupported "+caseType+" algorithm: %s\n", *algorithm)
		os.Exit(1)
	}

	// Use the shared RenameFiles function
	if err := rntocase.RenameFiles(files, converter, *dryRun, *interactive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
