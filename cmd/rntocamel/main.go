package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	strcase2 "github.com/ettle/strcase"
	"github.com/gobeam/stringy"
	"github.com/golang-cz/textcase"
	"github.com/iancoleman/strcase"
	strings2 "github.com/searKing/golang/go/strings"
	"github.com/tomeh/caseconv"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"strings"
)

const (
	defaultAlgo = "iancoleman"
	caseType    = "camel"
	appName     = "rntocamel"
)

var (
	wordSeparators []string
	algos          = map[string]func(string) (string, error){
		"iancoleman": func(s string) (string, error) {
			return strcase.ToCamel(s), nil
		},
		"lower-iancoleman": func(s string) (string, error) {
			return strcase.ToLowerCamel(s), nil
		},
		"lower-searking": func(s string) (string, error) {
			return strings2.LowerCamelCase(s), nil
		},
		"searking": func(s string) (string, error) {
			return strings2.UpperCamelCase(s), nil
		},
		"ettle": func(s string) (string, error) {
			return strcase2.ToCamel(s), nil
		},
		"go-ettle": func(s string) (string, error) {
			return strcase2.ToGoCamel(s), nil
		},
		"resenje": func(s string) (string, error) {
			return casbab.Camel(s), nil
		},
		"resenje-kebab": func(s string) (string, error) {
			return casbab.CamelKebab(s), nil
		},
		"resenje-snake": func(s string) (string, error) {
			return casbab.CamelSnake(s), nil
		},
		"tomeh": func(s string) (string, error) {
			return caseconv.ToCamel(s), nil
		},
		"golang-cz": func(s string) (string, error) {
			return textcase.CamelCase(s), nil
		},
		"gobeam": func(s string) (string, error) {
			return stringy.New(s).CamelCase(wordSeparators...).Get(), nil
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
	flag.Func("word-seperators", "Word separators. (gobeam only) Default: \"_-. \"", func(s string) error {
		wordSeparators = append(wordSeparators, s)
		return nil
	})
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
