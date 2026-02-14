# Bug Report: Issues with `gosubc` v0.0.17 for `rntocase`

## Summary

`gosubc` v0.0.17 cannot be used to upgrade `rntocase` fully because it lacks support for slice types (`[]string`, `[]int`, etc.) and custom flag parsing (`flag.Func`). This prevents upgrading tools like `rntocamel` that rely on these features (e.g., for `--word-seperators` or `--acronym`).

## Issues Found

### 1. Panic on Slice Types

`gosubc` panics when encountering a function parameter with a slice type, which is used in `rntocase` tools.

**Reproduction:**
Create a Go file with:
```go
// TestSlice is a subcommand `app slice`
// Flags:
//   names: --name -n
func TestSlice(names []string) {}
```
Run `gosubc generate`.

**Error:**
```
panic: Unsupported type: *ast.ArrayType

goroutine 1 [running]:
github.com/arran4/go-subcommand/parsers/commentv1.ParseGoFile(...)
    /home/jules/go/pkg/mod/github.com/arran4/go-subcommand@v0.0.17/parsers/commentv1/commentv1.go:283 +0x1bf4
...
```

### 2. Lack of Custom Flag Parsing (`flag.Func`)

`rntocase` uses `flag.Func` to handle complex flag logic or append to slices (e.g., `--word-seperators`). `gosubc` generates code using `flag.NewFlagSet` and adds flags based on function signature, but only supports basic types (`string`, `int`, `bool`, `time.Duration`). It doesn't appear to support custom types or `flag.Value` interface via comments.

Example from `cmd/rntocamel/main.go`:
```go
    flag.Func("word-seperators", "Word separators. (gobeam only) Default: \"_-. \"", func(s string) error {
        wordSeparators = append(wordSeparators, s)
        return nil
    })
```
Replacing this with `gosubc` would require dropping support for repeated flags or implementing custom flag parsing outside of `gosubc`'s control, which defeats the purpose of using `gosubc`.

## Changes in Interface

Even for tools that *could* be upgraded (like `rnreverse`), `gosubc` introduces interface changes:
- Help output format changes significantly.
- Adds `help`, `usage`, `version` subcommands automatically.
- Generated code structure (`cmd/<tool>/...`) is more complex than the current single-file structure.

## Conclusion

Upgrading to `gosubc` v0.0.17 is not feasible without either:
1.  Updating `gosubc` to support slice flags/`flag.Value`.
2.  Refactoring `rntocase` to remove features relying on complex flag parsing.

Therefore, `gosubc` v0.0.17 does not work for generating code for `rntocase` without major breaking features.
