package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	strcase2 "github.com/ettle/strcase"
	"github.com/iancoleman/strcase"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"strings"
)

const (
	defaultAlgo = "iancoleman"
	caseType    = "constant"
	appName     = "rntoconstant"
)

var (
	algos = map[string]func(string) (string, error){
		"iancoleman": func(s string) (string, error) {
			return strcase.ToScreamingSnake(s), nil
		},
		"ettle": func(s string) (string, error) {
			return strcase2.ToSNAKE(s), nil
		},
		"resenje": func(s string) (string, error) {
			return casbab.ScreamingSnake(s), nil
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
	keys := slices.Collect(maps.Keys(algos))
	slices.Sort(keys)
	algorithm := flag.String("algorithm", defaultAlgo, "Choose the "+caseType+" algorithm to use, supported: "+strings.Join(keys, ",")+".")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: "+appName+" [options] <file1> [<file2> ...]")
		_, _ = fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Conversion Examples:")
		rntocase.RenderUsageTable(algos)
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
		fmt.Printf("Unsupported "+caseType+" algorithm: %s\n", *algorithm)
		os.Exit(1)
	}

	// Use the shared RenameFiles function
	if err := rntocase.RenameFiles(files, converter, *dryRun, *interactive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
