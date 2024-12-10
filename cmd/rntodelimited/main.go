package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/gobeam/stringy"
	"github.com/iancoleman/strcase"
	strings2 "github.com/searKing/golang/go/strings"
	"maps"
	"os"
	"slices"
	"strings"
)

const (
	defaultAlgo      = "iancoleman"
	caseType         = "delimited"
	appName          = "rntodelimited"
	defaultIgnore    = ""
	defaultDelimiter = "."
)

func withDefault(def string, f func(s string) string) func(s string) string {
	return func(s string) string {
		if s == "" {
			s = def
		}
		return f(s)
	}
}

const (
	leadingString = "X"
)

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	flag.Func("acronym", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)", func(s string) error {
		strcase.ConfigureAcronym(s, s)
		return nil
	})
	flag.Func("acronym-from-file", "Words to consider acronyms and not to assume they are words, ie ID => ID rather than ID => Id (iancoleman only)", func(s string) error {
		return rntocase.LoadAcronymsFromFile(s)
	})
	ignore := flag.String("ignore", defaultIgnore, "Other delimiter characters to ignore when parsing")
	delimiter := flag.String("delimiter", defaultDelimiter, "Delimiter, default '.' but can be any single ascii character.")
	inputDelimiter := flag.String("input-delimiters", defaultDelimiter, "Input's Delimiters, default ' ' but can be multiple ascii character. (searking only)")
	var (
		wordSeparators []string
		algos          = map[string]func(string) (string, error){
			"iancoleman": func(s string) (string, error) {
				del := []byte(*delimiter)[0]
				return strcase.ToDelimited(s, del), nil
			},
			"screaming-iancoleman": func(s string) (string, error) {
				del := []byte(*delimiter)[0]
				return strcase.ToScreamingDelimited(s, del, *ignore, true), nil
			},
			"searking": func(s string) (string, error) {
				inDel := []rune(*inputDelimiter)
				return strings2.TransformCase(s, strings2.JoinGenerator(*delimiter, withDefault(leadingString, strings.ToLower)), inDel...), nil
			},
			"gobeam": func(s string) (string, error) {
				return stringy.New(s).Delimited(*delimiter, wordSeparators...).Get(), nil
			},
		}
	)
	flag.Func("word-seperators", "Word separators. (gobeam only) Default: \"_-. \"", func(s string) error {
		wordSeparators = append(wordSeparators, s)
		return nil
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
