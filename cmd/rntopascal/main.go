package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/ettle/strcase"
	"github.com/gobeam/stringy"
	"github.com/golang-cz/textcase"
	strings2 "github.com/searKing/golang/go/strings"
	"github.com/tomeh/caseconv"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"strings"
)

const (
	defaultAlgo = "ettle"
	caseType    = "pascal"
	appName     = "rntopascal"
)

var (
	wordSeparators []string
	algos          = map[string]func(string) (string, error){
		"searKing": func(s string) (string, error) {
			return strings2.PascalCase(s), nil
		},
		"ettle": func(s string) (string, error) {
			return strcase.ToPascal(s), nil
		},
		"go-ettle": func(s string) (string, error) {
			return strcase.ToGoPascal(s), nil
		},
		"resenje": func(s string) (string, error) {
			return casbab.Pascal(s), nil
		},
		"tomeh": func(s string) (string, error) {
			return caseconv.ToPascal(s), nil
		},
		"golang-cz": func(s string) (string, error) {
			return textcase.PascalCase(s), nil
		},
		"gobeam": func(s string) (string, error) {
			return stringy.New(s).PascalCase(wordSeparators...).Get(), nil
		},
	}
)

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	flag.Func("acronym", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id", func(s string) error {
		return nil
	})
	flag.Func("acronym-from-file", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id", func(s string) error {
		return rntocase.LoadAcronymsFromFile(s)
	})
	var wordSeparators []string
	flag.Func("word-seperators", "Word separators. (gobeam only) Default: \"_-. \"", func(s string) error {
		wordSeparators = append(wordSeparators, s)
		return nil
	})
	algorithm := flag.String("algorithm", defaultAlgo, "Choose the "+caseType+" algorithm to use, supported: "+strings.Join(slices.Collect(maps.Keys(algos)), ",")+".")
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
		fmt.Printf("Uunsupported "+caseType+" algorithm: %s\n", *algorithm)
		os.Exit(1)
	}

	// Use the shared RenameFiles function
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	if err := rntocase.RenameFiles(w, files, converter, *dryRun, *interactive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
