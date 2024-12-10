package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	strcase2 "github.com/ettle/strcase"
	"github.com/gobeam/stringy"
	"github.com/golang-cz/textcase"
	"github.com/iancoleman/strcase"
	"github.com/jedib0t/go-pretty/table"
	strings2 "github.com/searKing/golang/go/strings"
	"github.com/tomeh/caseconv"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"sort"
	"strings"
)

const (
	defaultAlgo = "iancoleman"
	caseType    = "kebab"
	appName     = "rntokebab"
)

var (
	wordSeparators []string
	algos          = map[string]func(string) (string, error){
		"iancoleman": func(s string) (string, error) {
			return strcase.ToKebab(s), nil
		},
		"searking": func(s string) (string, error) {
			return strings2.KebabCase(s), nil
		},
		"screaming-iancoleman": func(s string) (string, error) {
			return strcase.ToScreamingKebab(s), nil
		},
		"ettle": func(s string) (string, error) {
			return strcase2.ToKebab(s), nil
		},
		"go-ettle": func(s string) (string, error) {
			return strcase2.ToGoKebab(s), nil
		},
		"ETTLE": func(s string) (string, error) {
			return strcase2.ToKEBAB(s), nil
		},
		"resenje": func(s string) (string, error) {
			return casbab.Kebab(s), nil
		},
		"screaming-resenje": func(s string) (string, error) {
			return casbab.ScreamingKebab(s), nil
		},
		"resenje-camel": func(s string) (string, error) {
			return casbab.CamelKebab(s), nil
		},
		"resenje-snake": func(s string) (string, error) {
			return casbab.CamelSnake(s), nil
		},
		"tomeh": func(s string) (string, error) {
			return caseconv.ToKebab(s), nil
		},
		"golang-cz": func(s string) (string, error) {
			return textcase.KebabCase(s), nil
		},
		"gobeam": func(s string) (string, error) {
			return stringy.New(s).KebabCase(wordSeparators...).Get(), nil
		},
	}
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
		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"Algorithm"})
		keys := slices.Collect(maps.Keys(algos))
		sort.Strings(keys)
		for _, values := range rntocase.ExampleGroups {
			tw.AppendRow(slices.Collect(func(yield func(row any) bool) {
				yield("")
				if !yield(values.Name) {
					return
				}
				for _, value := range values.Values {
					if !yield(value) {
						return
					}
				}
			}))
			for _, key := range keys {
				tw.AppendRow(slices.Collect(func(yield func(row any) bool) {
					yield(key)
					yield("")
					for _, value := range values.Values {
						result, err := algos[key](value)
						if err != nil {
							panic(err)
						}
						if !yield(result) {
							return
						}
					}
				}))
			}
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
