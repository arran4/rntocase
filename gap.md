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

1. **Acronym Extraction (`rnacronym`)**: Converting "some nice string" to "sns". `strings2` focuses on case conversions rather than extracting initial letters. This feature is approximated using `strings2.Parse()` to extract words and pick the first letter.
2. **String Reversal (`rnreverse`)**: Reversing characters/runes in a string. `strings2` does not provide utility functions for reversing runes. We implement a custom rune-reversal function manually.
3. **Darwin Case (`rntodarwin`)**: Capitalizes the first letter of each word and separates words with `_` (e.g. `Some_String_Test`). `strings2` does not have an explicit Darwin case preset. This is approximated with `strings2.Parse()` and custom string joining.

### Behavioral Mismatches

As uncovered during testing by explicitly testing against original implementations (like `iancoleman/strcase`), `strings2` possesses subtle word boundary and dictionary parsing differences that produce different outputs in certain complex edge scenarios:
- **camelCase mismatch**: `strcase` outputs `HelloWorld` for leading underscores, while `strings2` outputs `helloWorld`. `strcase` outputs `HttpRequestId` for acronym handling while `strings2` outputs `httpRequestId`.
- **snake_case mismatch**: `strcase` parses `helloWorld` to `hello_world`, but `strings2` may retain capitalization on the second word when separating `hello_World`. It also handles double caps like `HTTPResponse` slightly differently (`http_response` vs `HTTP_Response`).
- **kebab-case mismatch**: Similar to snake-case, `helloWorld` translates to `hello-world` in `strcase` but to `hello-World` in `strings2`.
- **delimited mismatch**: `helloWorld` translates to `hello_world` via `strcase.ToDelimited`, while `strings2.ToFormattedString` translates to `hello_World`.
