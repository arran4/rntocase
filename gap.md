# Gap Analysis for Migration to `github.com/arran4/strings2`

The repository `rntocase` contains many individual executables that transform casing. Currently, they rely on various third-party casing libraries (like `strcase`, `stringy`, `caseconv`, `casbab`, etc.).
The goal is to replace these libraries with `github.com/arran4/strings2` where possible, while maintaining full support for existing case mutations, and to simplify the dependency tree.

## Gaps (Not Supported by `strings2`)

The following formats, features, flags or specific algorithm variations are dropped or require workarounds because they are not explicitly provided natively by `github.com/arran4/strings2`. Where `strings2` isn't a good fit, local implementations or the standard Go library are utilized.

### `strings2` Should NOT Implement (Out of Scope / Anti-Patterns)

1. **String Reversal (`rnreverse`)**: Reversing characters or words in a string.
   - **Local Implementation:** We implement custom rune-reversal functions manually and token reversal via `strings2.Parse()`.
   - **Recommendation:** Reversing strings is a geometric string translation rather than a formatting translation. Expanding `strings2` API surface area with reverse utilities falls outside typical casing package scope.

2. **Arbitrary Character Delimiter Overrides (`-input-delimiters`, `-word-seperators`)**: These were tied to specific libraries (`searking`, `gobeam`) and are no longer supported.
   - **Local Implementation:** We've dropped these specific CLI overrides to embrace the unified formatting model.
   - **Recommendation:** While mapping old arbitrary character delimiters to `strings2` configurations via custom partitioners is possible, it is brittle. `strings2` relies on robust token partitioner heuristics (camel boundaries, etc.). Forcing user-defined byte-level boundary arrays dilutes the strength of `strings2`'s AST tokenizer.
