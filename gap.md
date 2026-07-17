# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`. Where `strings2` isn't a good fit, local implementations or the standard Go library are utilized.

### `strings2` Should Implement (Feature Requests)

#### Feature Request #1: Simple Global Acronym Configuration

**The Gap:** Legacy tooling relied on `strcase` to define acronyms (e.g. `ID`, `HTTP`) from a file or CLI flags via a global configuration. While `strings2` supports acronyms inside `ParserConfig`, reading and injecting them easily from an external source dynamically via the API requires more manual implementation on the caller's side.

**Local Implementation:** We dropped CLI flag support for defining acronym files entirely to avoid writing boilerplate to map text inputs into `strings2.ParserConfig`.

**Implementation Suggestions for `strings2`:**
*Note: strings2 should manage its own configuration surface cleanly.*
- **Option A:** Provide an options factory `strings2.OptionAcronymList([]string{"HTTP", "JSON"})` that converts a slice of strings into `strings2` recognized acronym configurations internally mapping to the parser.
  ```go
  strings2.OptionAcronymList([]string{"HTTP", "JSON"})
  ```
- **Test Cases:**
  - Given loaded acronyms `["HTTP", "JSON"]`, input `Http Json Config` -> `HTTP_JSON_Config`

#### Feature Request #2: "Do Not Break" / Custom Ignore Patterns for Delimiters

**The Gap:** Users occasionally need explicit control over what `strings2` considers a word boundary or a delimiter to ignore. The removed `strcase` flags like `-ignore` mapped characters that should *not* cause a split (e.g. ignoring `.` so `foo.bar` stays `foo.bar` instead of becoming `foo_bar`).

**Implementation Suggestions for `strings2`:**
*Note: Providing an interface to define ignored boundaries is highly requested, but it must map cleanly to `strings2`'s parser without overriding its core AST logic.*
- **Option A:** Introduce `strings2.OptionIgnoreTokens(tokens ...string)` which tells the `strings2` partitioner to skip breaking on specific substrings or runes.
- **Option B:** Introduce a regex-based bypass `strings2.OptionRetainPattern(regexp.MustCompile(...))` which shields specific patterns from tokenization.
- **Test Cases:**
  - Given ignore token `.`: `my.namespace.config` -> `my.namespace.config` (instead of `my_namespace_config` in ToSnake)
  - Given ignore token `@`: `user@email` -> `user@email`

### `strings2` Should NOT Implement (Out of Scope / Anti-Patterns)

1. **String Reversal (`rnreverse`)**: Reversing characters or words in a string.
   - **Local Implementation:** We implement custom rune-reversal functions manually and token reversal via `strings2.Parse()`.
   - **Recommendation:** Reversing strings is a geometric string translation rather than a formatting translation. Expanding `strings2` API surface area with reverse utilities falls outside typical casing package scope.
