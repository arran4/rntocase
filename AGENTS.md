# AGENTS.md

## Repository Expectations for Agents

This document defines expectations for agents interacting with `arran4/rntocase`.
It emphasizes testing patterns, file system abstractions, and integration tests.

### Filesystem Architecture

- Filesystem-heavy business logic should not directly depend on host `os.*` side effects when a small injectable capability interface is practical.
- Keep OS calls in explicit adapter/boundary layers.
- Use standard `io/fs` interfaces where sufficient; add only narrowly scoped writable/mutation interfaces when necessary.
- Preserve existing call sites where practical by using optional dependency injection via `ops ...any` rather than adding mandatory plumbing solely for tests.
- Do not create oversized generic filesystem abstractions. Define exactly what your code needs.

### Testing

- Memory filesystem tests are the default for filesystem business logic.
- Prefer `fstest.MapFS` for read-only scenarios.
- Use a small repository-local mock/writable FS for mutations when necessary.
- Prefer embedded `txtar` fixtures for non-trivial directory trees or input/expected filesystem scenarios.
- Use deterministic fixture walking and one `t.Run` per `txtar` case.
- `t.TempDir()` is appropriate for genuine integration/OS-semantics tests, not simply because production currently calls `os.*`.
- Keep a small real-filesystem integration suite to verify OS adapters, rename behavior, and symlink semantics.
- New features that introduce filesystem traversal or mutation should normally demonstrate both memory-level behavioral coverage and targeted OS-bound integration coverage when relevant.

### Agent review checklist

Before considering filesystem-related work complete, agents should ask:

1. Is business logic unnecessarily calling `os.*` directly?
2. Could the behavior be tested through `fs.FS` or a small injected interface?
3. Are basic tests creating real files unnecessarily?
4. Would a txtar fixture make a complex filesystem scenario easier to review?
5. Are real disk tests limited to behavior that genuinely requires OS semantics?
6. Are traversal and fixture results deterministic?
7. Did the change preserve production API compatibility where practical?

### Reference Materials

For further context and practical examples, agents should refer to these articles:

- [Go Memory FSs Everywhere in Test: Optional Dependency Injection via Type-Switched Variadic Args](https://arran4.github.io/blog/post/2026/020-optional-dependency-injection-via-type-switched-variadic-args/)
- [Testing File Systems: How I use MockFS, MapFS, and SimpleFS in Go](https://arran4.github.io/blog/post/2026/033-testing-fs-with-mapfs-mockfs/)
- [Txtar Patterns for Agents: Data, Scenarios, and Embedded Walkers](https://arran4.github.io/blog/post/2026/004-txtar-patterns-for-agents/)

### Limitations

- Memory MockFS test fixtures use UNIX-style path mappings (e.g., `/testdir`). Platform-dependent logic inside `filepath.Join` could theoretically differ on Windows, so the core `MockFS` limits full cross-platform parity. Rely on OS-integration specs (`t.TempDir()`) for strict cross-platform validations.
