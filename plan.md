1. **Fix atomic updates:** Update `updateSingleSkill` in `cmd/rntocase/skill.go` to download the tarball before removing the local directory.
2. **Fix nested skill listing:** Update `ListInstalledSkills` in `internal/skill/manager.go` to use `filepath.Walk` (or `fs.WalkDir`) to find `.rntocase-skill.json` files recursively rather than just depth=1.
3. **Complete pre commit steps.**
