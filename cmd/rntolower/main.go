package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/gobeam/stringy"
	"github.com/jedib0t/go-pretty/table"
	skstrings "github.com/searKing/golang/go/strings"
	"maps"
	"os"
	"resenje.org/casbab"
	"slices"
	"sort"
	"strings"
)

const (
	defaultAlgo = "resenje"
	caseType    = "lowercase"
	appName     = "rntolower"
)

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	var (
		algos = map[string]func(string) (string, error){
			"searking": func(s string) (string, error) {
				return skstrings.LowerCase(s), nil
			},
			"go": func(s string) (string, error) {
				return strings.ToLower(s), nil
			},
			"resenje": func(s string) (string, error) {
				return casbab.Lower(s), nil
			},
			"gobeam": func(s string) (string, error) {
				return stringy.New(s).ToLower(), nil
			},
		}
	)

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
