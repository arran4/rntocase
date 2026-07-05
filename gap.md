# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Migrated / Supported Conversions

The following commands and cases are natively supported by `github.com/arran4/strings2` and standard library and will be migrated to `strings2`:
- `rntocamel` (using `strings2.ToPascal` to match default outputs of old library)
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

1. **Acronym Extraction (`rnacronym`)**: Converting "some nice string" to "sns". `strings2` focuses on case conversions rather than extracting initial letters. This feature is approximated using `strings2.Parse()` to extract words and pick the first letter.
2. **String Reversal (`rnreverse`)**: Reversing characters/runes in a string. `strings2` does not provide utility functions for reversing runes. We implement a custom rune-reversal function manually.
3. **Darwin Case (`rntodarwin`)**: Capitalizes the first letter of each word and separates words with `_` (e.g. `Some_String_Test`). `strings2` does not have an explicit Darwin case preset. This is approximated with `strings2.Parse()` and custom string joining.
4. **Acronym flags (`-acronym`, `-acronym-from-file`)**: Several tools exposed these flags from `iancoleman/strcase`. Since we dropped `strcase`, these flags are no longer supported and have been removed. `strings2` does support acronyms via options but integrating it globally as a flag replacement is a gap.
5. **Other specific flags (`-ignore`, `-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`iancoleman`, `searking`, `gobeam`) and are no longer supported.
6. **Lower Camel Case (`rntocamel`)**: Explicit lower-camel variant (from `strcase.ToLowerCamel`) was dropped to unify under `strings2.ToCamel`.
7. **Screaming Kebab / Screaming Delimited / Lowercase Snake (`rntokebab`, `rntodelimited`, `rntosnake`)**: Specific case variations provided by dropped libraries (e.g. `screaming-kebab`) were removed to simplify the codebase to the primary case type for each tool.

### Behavioral Mismatches

With `strings2@v0.0.6`, many behavioral differences relative to `strcase` were mitigated by exposing `CMWhispering` and other configuration modes. The tools have been configured to match `strcase` where possible, eliminating the prior known output mismatch skips except where the underlying parser intrinsically defines acronyms or punctuation slightly differently from `strcase` (which are noted as skipped tests in `cmd/*_test.go` suites where applicable).
