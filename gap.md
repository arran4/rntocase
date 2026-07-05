# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`. Where `strings2` isn't a good fit, local implementations or the standard Go library are utilized.

### `strings2` Should Implement (Feature Requests)

#### Feature Request #1: Native Darwin Case (`ToDarwinCase`)

**The Gap:** Darwin case capitalizes the first letter of every word and separates them with an underscore. E.g. `Some_String_Test`.

**Local Implementation:** We achieve this cleanly by composing options natively:
```go
strings2.ToFormattedString(s, strings2.OptionDelimiter("_"), strings2.OptionCaseMode(strings2.CMAllTitle))
```

**Implementation Suggestions for `strings2`:**
*Note: These are only suggestions. The implementation should align with `strings2`'s internal style, provided the core need is addressed.*
- **Option A:** Add a direct alias function like `ToSnake` and `ToCamel` do.
  ```go
  func ToDarwinCase(input string, opts ...any) (string, error) {
      opts = append(opts, OptionDelimiter("_"), OptionCaseMode(CMAllTitle))
      return ToFormattedString(input, opts...)
  }
  ```
- **Test Cases:**
  - `hello world` -> `Hello_World`
  - `camelCaseInput` -> `Camel_Case_Input`
  - `mixed-UP-Kebab` -> `Mixed_Up_Kebab`

#### Feature Request #2: Smart Title Case (`ToTitleCase`)

**The Gap:** `strings2` lacks a smart title-casing mode. True title-casing capitalizes the first letter of words but typically leaves minor conjunctions and prepositions lowercase (e.g. "a", "an", "the", "in", "of").

**Local Implementation:** We utilize the standard Go library (`golang.org/x/text/cases`) to handle title-casing strings as a fallback, completely bypassing `strings2`.

**Implementation Suggestions for `strings2`:**
*Note: These are only suggestions. The implementation should align with `strings2`'s internal style, provided the core need is addressed.*
- **Option A:** Introduce `CMTitle` or `CMSmartTitle` to `CaseMode` which evaluates word context during formatting.
  ```go
  // Internal formatter checks if word matches a skip-list dictionary
  func ToTitleCase(input string, opts ...any) (string, error) {
      // ...
  }
  ```
- **Option B:** Supply an `OptionSkipTitleWords([]string{"in", "the"})` configuration to allow consumers to define what words stay lowercase during an `OptionCaseMode(CMAllTitle)` conversion.
- **Test Cases:**
  - `the lord of the rings` -> `The Lord of the Rings`
  - `A_NEW_HOPE` -> `A New Hope`

#### Feature Request #3: Simple Global/File Acronym Configuration

**The Gap:** Legacy tooling relied on `strcase` to define acronyms (e.g. `ID`, `HTTP`) from a file or CLI flags via a global configuration. While `strings2` supports acronyms inside `ParserConfig`, reading and injecting them easily from an external source dynamically via the API requires more manual implementation on the caller's side.

**Local Implementation:** We dropped CLI flag support for defining acronym files entirely to avoid writing boilerplate to map text inputs into `strings2.ParserConfig`.

**Implementation Suggestions for `strings2`:**
*Note: These are only suggestions. The implementation should align with `strings2`'s internal style, provided the core need is addressed.*
- **Option A:** Provide an options helper that converts a slice of strings into `strings2` recognized acronym configurations directly.
  ```go
  strings2.OptionAcronymList([]string{"HTTP", "JSON"})
  ```
- **Option B:** Provide a utility function `LoadAcronymsFromFile(filepath)` inside the module that handles file io mapping to tokens seamlessly.
- **Test Cases:**
  - Given loaded acronyms `["HTTP", "JSON"]`, input `Http Json Config` -> `HTTP_JSON_Config`


### `strings2` Should NOT Implement (Out of Scope / Anti-Patterns)

1. **String Reversal (`rnreverse`)**: Reversing characters or words in a string.
   - **Local Implementation:** We implement custom rune-reversal functions manually and token reversal via `strings2.Parse()`.
   - **Recommendation:** Reversing strings is a geometric string translation rather than a formatting translation. Expanding `strings2` API surface area with reverse utilities falls outside typical casing package scope.

2. **Arbitrary Character Delimiter Overrides (`-ignore`, `-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`iancoleman`, `searking`, `gobeam`) and are no longer supported.
   - **Local Implementation:** We've dropped these specific CLI overrides to embrace the unified formatting model.
   - **Recommendation:** While mapping old arbitrary character delimiters to `strings2` configurations via custom partitioners is possible, it is brittle. `strings2` relies on robust token partitioner heuristics (camel boundaries, etc.). Forcing user-defined byte-level boundary arrays dilutes the strength of `strings2`'s AST tokenizer.
