package cli

import (
	"github.com/arran4/rntocase"
	"github.com/arran4/strings2"
	"strings"
)




// RunAcronym is a subcommand `rntocase acronym` -- Rename files by acronym
//
func RunAcronym(dryRun bool, interactive bool, files ...string) error {
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
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunCamel is a subcommand `rntocase camel` -- Rename files to camel case
//
func RunCamel(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToPascal(s, strings2.ParserEmitEmpty(true))
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunConstant is a subcommand `rntocase constant` -- Rename files to constant case
//
func RunConstant(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		val, err := strings2.ToSnake(s, strings2.ParserEmitEmpty(true))
		if err != nil {
			return "", err
		}
		return strings.ToUpper(val), nil
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunDarwin is a subcommand `rntocase darwin` -- Rename files to darwin case
//
func RunDarwin(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToDarwin(s)
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunDelimited is a subcommand `rntocase delimited` -- Rename files with a custom delimiter
//
func RunDelimited(delimiter string, ignore string, dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToFormattedString(s, strings2.OptionDelimiter(delimiter), strings2.OptionIgnore(ignore), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}

	// Re-assign usage after parse to ensure variables are evaluated correctly during -h
	// We have to explicitly handle usage if requested since Parse is already called. We'll ignore the double usage logic since Parse stops on error.

	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunDot is a subcommand `rntocase dot` -- Rename files to dot case
//
func RunDot(delimiter string, dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToFormattedString(s, strings2.OptionDelimiter(delimiter), strings2.OptionFirstLower(), strings2.ParserEmitEmpty(true))
	}

	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunKebab is a subcommand `rntocase kebab` -- Rename files to kebab case
//
func RunKebab(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToKebab(s, strings2.OptionLoose(), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunLower is a subcommand `rntocase lower` -- Rename files to lower case
//
func RunLower(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings.ToLower(s), nil
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunLowerLeading is a subcommand `rntocase lowerleading` -- Rename files with a lower leading character
//
func RunLowerLeading(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.LowerCaseFirstWithErr(s)
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunPascal is a subcommand `rntocase pascal` -- Rename files to pascal case
//
func RunPascal(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToPascal(s)
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunReverse is a subcommand `rntocase reverse` -- Reverse characters or words in file names
//
func RunReverse(wordMode bool, dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	}



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
	if wordMode {
		activeConverter = wordReverseConverter
	}

	return rntocase.RenameFiles(files, activeConverter, dryRun, interactive)
}

// RunSnake is a subcommand `rntocase snake` -- Rename files to snake case
//
func RunSnake(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToSnake(s, strings2.OptionLoose(), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunTitle is a subcommand `rntocase title` -- Rename files to title case
//
func RunTitle(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.ToTitle(s)
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunUpper is a subcommand `rntocase upper` -- Rename files to upper case
//
func RunUpper(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings.ToUpper(s), nil
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunUpperLeading is a subcommand `rntocase upperleading` -- Rename files with an upper leading character
//
func RunUpperLeading(dryRun bool, interactive bool, files ...string) error {
	converter := func(s string) (string, error) {
		return strings2.UpperCaseFirstWithErr(s)
	}
	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}

// RunTrim is a subcommand `rntocase trim` -- Trim whitespace or specific characters from file names
//
func RunTrim(algorithm string, trimChars string, dryRun bool, interactive bool, files ...string) error {


	converter := func(s string) (string, error) {
		if trimChars == "" {
			return strings.TrimSpace(s), nil
		} else {
			return strings.Trim(s, trimChars), nil
		}
	}

	return rntocase.RenameFiles(files, converter, dryRun, interactive)
}
