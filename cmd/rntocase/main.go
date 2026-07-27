package main

import (
	"flag"
	"fmt"
	"github.com/arran4/rntocase"
	"github.com/arran4/strings2"
	"os"
	"strings"
)

const appName = "rntocase"

type subcommand struct {
	name        string
	description string
	run         func(args []string) error
}

func main() {
	cmds := []subcommand{
		{"acronym", "Rename files by acronym", runAcronym},
		{"camel", "Rename files to camel case", runCamel},
		{"constant", "Rename files to constant case", runConstant},
		{"darwin", "Rename files to darwin case", runDarwin},
		{"delimited", "Rename files with a custom delimiter", runDelimited},
		{"dot", "Rename files to dot case", runDot},
		{"kebab", "Rename files to kebab case", runKebab},
		{"lower", "Rename files to lower case", runLower},
		{"lowerleading", "Rename files with a lower leading character", runLowerLeading},
		{"pascal", "Rename files to pascal case", runPascal},
		{"reverse", "Reverse characters or words in file names", runReverse},
		{"snake", "Rename files to snake case", runSnake},
		{"title", "Rename files to title case", runTitle},
		{"upper", "Rename files to upper case", runUpper},
		{"upperleading", "Rename files with an upper leading character", runUpperLeading},
		{"trim", "Trim whitespace or specific characters from file names", runTrim},
		{"skill", "Manage AI agent skills for this CLI", runSkill},
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> [options] <file1> [<file2> ...]\n\n", appName)
		fmt.Fprintf(os.Stderr, "Subcommands:\n")
		for _, cmd := range cmds {
			fmt.Fprintf(os.Stderr, "  %-15s %s\n", cmd.name, cmd.description)
		}
		fmt.Fprintf(os.Stderr, "\nRun '%s <subcommand> -h' for more details on a specific subcommand.\n", appName)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	subcommandName := os.Args[1]

	// Handle global help
	if subcommandName == "-h" || subcommandName == "--help" || subcommandName == "help" {
		flag.Usage()
		os.Exit(0)
	}

	for _, cmd := range cmds {
		if cmd.name == subcommandName {
			if err := cmd.run(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcommandName)
	flag.Usage()
	os.Exit(1)
}

func setupFlags(cmdName string, args []string, extraFlags func(*flag.FlagSet), converter func(string) (string, error)) (*flag.FlagSet, *bool, *bool) {
	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, "Display the intended changes without renaming.")
	interactive := fs.Bool("interactive", false, "Ask for confirmation before renaming each file.")
	if extraFlags != nil {
		extraFlags(fs)
	}

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [options] <file1> [<file2> ...]\n\nOptions:\n", appName, cmdName)
		fs.PrintDefaults()
		if converter != nil {
			fmt.Fprintf(os.Stderr, "\nConversion Examples:\n")
			rntocase.RenderUsageTable(converter)
		}
	}

	_ = fs.Parse(args)

	return fs, dryRun, interactive
}

func validateFiles(fs *flag.FlagSet) []string {
	files := fs.Args()
	if len(files) == 0 {
		fmt.Println("Error: No files provided.")
		fs.Usage()
		os.Exit(1)
	}
	return files
}

func runAcronym(args []string) error {
	converter := func(s string) (string, error) {
		words, err := strings2.Parse(s)
		if err != nil {
			return "", err
		}
		var result strings.Builder
		for _, w := range words {
			for _, r := range w.String() {
				result.WriteString(strings.ToUpper(string(r)))
				break
			}
		}
		return result.String(), nil
	}

	fs, dryRun, interactive := setupFlags("acronym", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runCamel(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToPascal(s, strings2.ParserEmitEmpty(true))
	}
	fs, dryRun, interactive := setupFlags("camel", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runConstant(args []string) error {
	converter := func(s string) (string, error) {
		val, err := strings2.ToSnake(s, strings2.ParserEmitEmpty(true))
		if err != nil {
			return "", err
		}
		return strings.ToUpper(val), nil
	}
	fs, dryRun, interactive := setupFlags("constant", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runDarwin(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToDarwin(s)
	}
	fs, dryRun, interactive := setupFlags("darwin", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runDelimited(args []string) error {
	var delimiter *string
	var ignore *string

	// Create a dummy converter to populate the table correctly since flags aren't parsed yet.
	// But it will panic/error out or just show _ since we rely on flag pointers. We can make a localized version.
	var converter func(string) (string, error)
	fs, dryRun, interactive := setupFlags("delimited", args, func(fs *flag.FlagSet) {
		delimiter = fs.String("delimiter", "_", "The delimiter string to separate words")
		ignore = fs.String("ignore", "", "Characters to ignore when breaking boundaries")
	}, nil)

	converter = func(s string) (string, error) {
		return strings2.ToFormattedString(s, strings2.OptionDelimiter(*delimiter), strings2.OptionIgnore(*ignore), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}

	// Re-assign usage after parse to ensure variables are evaluated correctly during -h
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [options] <file1> [<file2> ...]\n\nOptions:\n", appName, "delimited")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nConversion Examples:\n")
		rntocase.RenderUsageTable(converter)
	}
	// We have to explicitly handle usage if requested since Parse is already called. We'll ignore the double usage logic since Parse stops on error.

	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runDot(args []string) error {
	var delimiter *string

	var converter func(string) (string, error)
	fs, dryRun, interactive := setupFlags("dot", args, func(fs *flag.FlagSet) {
		delimiter = fs.String("delimiter", ".", "The delimiter string to separate words")
	}, nil)

	converter = func(s string) (string, error) {
		return strings2.ToFormattedString(s, strings2.OptionDelimiter(*delimiter), strings2.OptionFirstLower(), strings2.ParserEmitEmpty(true))
	}

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [options] <file1> [<file2> ...]\n\nOptions:\n", appName, "dot")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nConversion Examples:\n")
		rntocase.RenderUsageTable(converter)
	}

	files := validateFiles(fs)
	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runKebab(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToKebab(s, strings2.OptionLoose(), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}
	fs, dryRun, interactive := setupFlags("kebab", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runLower(args []string) error {
	converter := func(s string) (string, error) {
		return strings.ToLower(s), nil
	}
	fs, dryRun, interactive := setupFlags("lower", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runLowerLeading(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.LowerCaseFirstWithErr(s)
	}
	fs, dryRun, interactive := setupFlags("lowerleading", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runPascal(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToPascal(s)
	}
	fs, dryRun, interactive := setupFlags("pascal", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runReverse(args []string) error {
	var wordMode *bool

	converter := func(s string) (string, error) {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	}

	fs, dryRun, interactive := setupFlags("reverse", args, func(fs *flag.FlagSet) {
		wordMode = fs.Bool("words", false, "Reverse the order of words instead of the characters.")
	}, converter)

	wordReverseConverter := func(s string) (string, error) {
		words, err := strings2.Parse(s)
		if err != nil {
			return "", err
		}
		for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
			words[i], words[j] = words[j], words[i]
		}
		return strings2.WordsToFormattedCase(words, strings2.OptionCaseMode(strings2.CMVerbatim), strings2.OptionDelimiter(" "))
	}

	activeConverter := converter
	if *wordMode {
		activeConverter = wordReverseConverter
	}

	files := validateFiles(fs)
	return rntocase.RenameFiles(files, activeConverter, *dryRun, *interactive)
}

func runSnake(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToSnake(s, strings2.OptionLoose(), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}
	fs, dryRun, interactive := setupFlags("snake", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runTitle(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.ToTitle(s)
	}
	fs, dryRun, interactive := setupFlags("title", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runUpper(args []string) error {
	converter := func(s string) (string, error) {
		return strings.ToUpper(s), nil
	}
	fs, dryRun, interactive := setupFlags("upper", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runUpperLeading(args []string) error {
	converter := func(s string) (string, error) {
		return strings2.UpperCaseFirstWithErr(s)
	}
	fs, dryRun, interactive := setupFlags("upperleading", args, nil, converter)
	files := validateFiles(fs)

	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}

func runTrim(args []string) error {
	var trimChars *string

	fs, dryRun, interactive := setupFlags("trim", args, func(fs *flag.FlagSet) {
		fs.String("algorithm", "go", "Choose the trim algorithm to use, supported: go.")
		trimChars = fs.String("trim", "", "Characters to trim off end and start of name white space if not set.")
	}, nil)

	converter := func(s string) (string, error) {
		if *trimChars == "" {
			return strings.TrimSpace(s), nil
		} else {
			return strings.Trim(s, *trimChars), nil
		}
	}

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [options] <file1> [<file2> ...]\n\nOptions:\n", appName, "trim")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nConversion Examples:\n")
		rntocase.RenderUsageTable(converter)
	}

	files := validateFiles(fs)
	return rntocase.RenameFiles(files, converter, *dryRun, *interactive)
}
