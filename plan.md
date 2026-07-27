1. **Fix `errcheck` violations**: Add explicit error ignores (`_ = ...`) to defer statements and close operations in `hash.go`, `installer.go`, `installer_test.go`, `manager.go`, `manifest_test.go`.
2. **Fix `staticcheck` violations**:
   - `SA4006`: In `cmd/rntocase/skill.go:248`, `sha` is unused because `updateSingleSkill` ignores it initially.
   - `QF1003`: Refactor `internal/skill/target.go` to use a tagged `switch scope { ... }` block instead of an `if-else` chain.
3. **Submit the fixes and reply to the PR comment**.
