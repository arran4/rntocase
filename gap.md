# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`. Where `strings2` isn't a good fit, local implementations or the standard Go library are utilized.

1. **Acronym Extraction (`rnacronym`)**: Converting "some nice string" to "sns". `strings2` focuses on case conversions rather than extracting initial letters.
   - **Local Implementation:** This feature is approximated using `strings2.Parse()` to extract words and pick the first letter.
   - **Feature Request:** While `strings2` could technically add an `AcronymCase` mode or a specific `OptionAcronymExtraction`, it's not highly advisable as it strays from standard delimited casing semantics. Our current explicit loop is sufficient and clear.

2. **Darwin Case (`rntodarwin`)**: Capitalizes the first letter of each word and separates words with `_` (e.g. `Some_String_Test`). `strings2` does not have an explicit Darwin case preset.
   - **Local Implementation:** This is natively approximated in `rntodarwin` with `strings2.ToFormattedString(s, strings2.OptionDelimiter("_"), strings2.OptionCaseMode(strings2.CMAllTitle))`.
   - **Feature Request:** `strings2` could easily support this natively by adding a `ToDarwinCase` function wrapping those exact options as an alias.

3. **Title Case (`rntotitle`)**: Standard title-casing.
   - **Local Implementation:** `strings2` does not strictly map to smart title case conventions natively. A basic standard-library equivalent (`cases.Title`) from `golang.org/x/text` is utilized instead as a robust fallback.
   - **Feature Request:** Adding smart title casing (handling prepositions correctly) to `strings2` could be highly valuable.

4. **Acronym flags (`-acronym`, `-acronym-from-file`)**: Several tools exposed these flags from `iancoleman/strcase`.
   - **Local Implementation:** Since we dropped `strcase`, these flags are no longer supported and have been removed.
   - **Feature Request:** `strings2` could support this via its smart acronyms config, but parsing files manually and injecting them into `strings2.ParserConfig` arrays might overcomplicate the standard mapping layer.

5. **Other specific flags (`-ignore`, `-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`iancoleman`, `searking`, `gobeam`) and are no longer supported.
   - **Feature Request:** `strings2` has its own robust token partitioner config; mapping old arbitrary character delimiters to `strings2` configurations is possible but potentially fragile. We've dropped them to enforce `strings2`'s native heuristics.
