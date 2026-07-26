package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/arran4/strings2"
	"os"
)

const (
	appName     = "rnreverse"
)

func converter(s string) (string, error) {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func wordReverseConverter(s string) (string, error) {
    words, err := strings2.Parse(s)
    if err != nil {
        return "", err
    }
    for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
        words[i], words[j] = words[j], words[i]
    }
    return strings2.WordsToFormattedCase(words, strings2.OptionCaseMode(strings2.CMVerbatim), strings2.OptionDelimiter(" "))
}

func main() {
	// Define flags
	dryRun := flag.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := flag.Bool("interactive", false, "Ask for confirmation before renaming each file.")
    wordMode := flag.Bool("words", false, "Reverse the order of words instead of the characters.")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: "+appName+" [options] <file1> [<file2> ...]")
		_, _ = fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr)
		_, _ = fmt.Fprintln(os.Stderr, "Conversion Examples:")
		rntocase.RenderUsageTable(converter)
	}
	flag.Parse()

	// Get file arguments
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Error: No files provided.")
		flag.Usage()
		os.Exit(1)
	}

    activeConverter := converter
    if *wordMode {
        activeConverter = wordReverseConverter
    }

	// Use the shared RenameFiles function
	if err := rntocase.RenameFiles(files, activeConverter, *dryRun, *interactive); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
