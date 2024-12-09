package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/iancoleman/strcase"
	strings2 "github.com/searKing/golang/go/strings"
	"maps"
	"os"
	"slices"
	"strings"
)

const (
	defaultAlgo = "iancoleman"
	caseType    = "kebab"
	appName     = "rntokebab"
)

var (
	algos = map[string]func(string) (string, error){
		"iancoleman": func(s string) (string, error) {
			return strcase.ToKebab(s), nil
		},
		"searking": func(s string) (string, error) {
			return strings2.KebabCase(s), nil
		},
		"screaming-iancoleman": func(s string) (string, error) {
			return strcase.ToScreamingKebab(s), nil
		},
	}
)

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	flag.Func("acronym", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id", func(s string) error {
		strcase.ConfigureAcronym(s, s)
		return nil
	})
	flag.Func("acronym-from-file", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id", func(s string) error {
		return rntocase.LoadAcronymsFromFile(s)
	})
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
