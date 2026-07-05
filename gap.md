# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`. Where `strings2` isn't a good fit, local implementations or the standard Go library are utilized.

### `strings2` Should Implement (Feature Requests)

1. **Darwin Case (`rntodarwin`)**: Capitalizes the first letter of each word and separates words with `_` (e.g. `Some_String_Test`). `strings2` does not have an explicit Darwin case preset.
   - **Local Implementation:** This is natively approximated in `rntodarwin` with `strings2.ToFormattedString(s, strings2.OptionDelimiter("_"), strings2.OptionCaseMode(strings2.CMAllTitle))`.
   - **Recommendation:** `strings2` could easily support this natively by adding a `ToDarwinCase` function wrapping those exact options as a clean alias.

2. **Title Case (`rntotitle`)**: Standard title-casing.
   - **Local Implementation:** `strings2` does not strictly map to smart title case conventions natively. A basic standard-library equivalent (`cases.Title`) from `golang.org/x/text` is utilized instead as a robust fallback.
   - **Recommendation:** Adding smart title casing (handling prepositions correctly like "in", "the", etc.) to `strings2` could be highly valuable for users building templating generators.

3. **Global Acronym Configuration (`-acronym-from-file`)**:
   - **Local Implementation:** Since we dropped `strcase`, dynamic CLI acronym file parsing flags are no longer supported and have been removed.
   - **Recommendation:** `strings2` could implement a cleaner global acronym ingestion API or CLI tooling parser to make managing `strings2.ParserConfig`'s smart acronyms less complex for implementers.

### `strings2` Should NOT Implement (Out of Scope / Anti-Patterns)

1. **Acronym Extraction (`rnacronym`)**: Converting "some nice string" to "sns". `strings2` focuses on case conversions rather than extracting initial letters.
   - **Local Implementation:** This feature is cleanly approximated using `strings2.Parse()` to extract words and pick the first letter.
   - **Recommendation:** While `strings2` could technically add an `AcronymCase` mode or an `OptionAcronymExtraction`, it strays too far from standard delimited casing semantics. Explicit loops are sufficient.

2. **String Reversal (`rnreverse`)**: Reversing characters or words in a string.
   - **Local Implementation:** We implement custom rune-reversal functions manually and token reversal via `strings2.Parse()`.
   - **Recommendation:** Reversing strings is a geometric string translation rather than a formatting translation. Expanding `strings2` API surface area with reverse utilities falls outside typical casing package scope.

3. **Arbitrary Character Delimiter Overrides (`-ignore`, `-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`iancoleman`, `searking`, `gobeam`) and are no longer supported.
   - **Local Implementation:** We've dropped these specific CLI overrides to embrace the unified formatting model.
   - **Recommendation:** While mapping old arbitrary character delimiters to `strings2` configurations via custom partitioners is possible, it is brittle. `strings2` relies on robust token partitioner heuristics (camel boundaries, etc.). Forcing user-defined byte-level boundary arrays dilutes the strength of `strings2`'s AST tokenizer.
