# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2@v0.0.4` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Migrated / Supported Conversions

The following commands and cases are natively supported by `github.com/arran4/strings2` and standard library and will be migrated to `strings2`:
- `rntocamel` (using `strings2.ToCamel`)
- `rntoconstant` (using `strings2.ToSnake` + `strings.ToUpper` for screaming snake case)
- `rntodelimited` (using `strings2.ToFormattedString(s, strings2.OptionDelimiter(...))`)
- `rntodot` (using `strings2.ToFormattedString(s, strings2.OptionDelimiter("."))`)
- `rntokebab` (using `strings2.ToKebab`)
- `rntopascal` (using `strings2.ToPascal`)
- `rntosnake` (using `strings2.ToSnake`)
- `rntolowerleading` (using `strings2.LowerCaseFirst`)
- `rntoupperleading` (using `strings2.UpperCaseFirst`)
- `rntolower`, `rntoupper`, `rntrim` (using standard library `strings` package)

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`:

1. **Title Case (`rntotitle`)**: Standard title-casing. `strings.Title` is deprecated, and `strings.ToTitle` maps to uppercase. `cases.Title` doesn't fully handle "smart" titlecasing rules natively. `strings2` does not strictly map to smart title case conventions natively. A basic standard-library equivalent (`cases.Title`) or custom parsing will be necessary.
2. **Acronym Extraction (`rnacronym`)**: Converting "some nice string" to "sns". `strings2` focuses on case conversions rather than extracting initial letters. This feature is approximated using `strings2.Parse()` to extract words and pick the first letter.
3. **String Reversal (`rnreverse`)**: Reversing characters/runes in a string. `strings2` does not provide utility functions for reversing runes. We implement a custom rune-reversal function manually.
4. **Darwin Case (`rntodarwin`)**: Capitalizes the first letter of each word and separates words with `_` (e.g. `Some_String_Test`). `strings2` does not have an explicit Darwin case preset. This is approximated with `strings2.Parse()` and custom string joining.
5. **Acronym flags (`-acronym`, `-acronym-from-file`)**: Several tools exposed these flags from `iancoleman/strcase`. Since we dropped `strcase`, these flags are no longer supported and have been removed. `strings2` does support acronyms via options but integrating it globally as a flag replacement is a gap.
6. **Other specific flags (`-ignore`, `-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`iancoleman`, `searking`, `gobeam`) and are no longer supported.
7. **Lower Camel Case (`rntocamel`)**: Explicit lower-camel variant (from `strcase.ToLowerCamel`) was dropped to unify under `strings2.ToCamel`.
8. **Screaming Kebab / Screaming Delimited / Lowercase Snake (`rntokebab`, `rntodelimited`, `rntosnake`)**: Specific case variations provided by dropped libraries (e.g. `screaming-kebab`) were removed from the `algos` map to simplify the codebase to the primary case type for each tool.
