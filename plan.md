1.  **Refactor `route` job:**
    *   Change defaults to `false` for everything except where explicit logic dictates it.
    *   For `workflow_dispatch` with `publish-tag`, explicitly set `run_code_checks=true`, `run_build=true`, `run_publisher=true`.
    *   Set `workflow_dispatch.inputs.mode.required: true` and add a default mode `build` (which should run code checks and build).
    *   Ensure `lint-fix` routes strictly to autofix and not regular code checks.
    *   Drop `reopened` / `ready_for_review` / `closed` triggers if the `pull_request` event has them.
2.  **Refactor `release-ready` gate:**
    *   Require explicit `success` states for `validation` and `build`. `skipped` is not acceptable.
    *   Ensure `release-ready` only runs if `run_release == 'true'` or `run_publisher == 'true'`.
3.  **Fix `GITHUB_*` variables:**
    *   Stop overriding `GITHUB_REF_NAME` and `GITHUB_SHA`. Use variables like `AUTHORITATIVE_BRANCH` and fetch against that.
4.  **Concurrency controls:**
    *   Add a non-canceling serialization group to `prepare-release-tag`: `group: ${{ github.workflow }}-release-preparation`.
    *   Add a non-canceling serialization group to `publisher`: `group: ${{ github.workflow }}-publish-${{ github.ref }}`.
5.  **Regex fix:**
    *   Fix explicit version validation regex in bash script: `^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$`.
6.  **Dependency Versions:**
    *   Switch `goreleaser` action to use `version: '~> v2'`.
    *   Switch `actions/upload-artifact` to `@v4` (it is already `@v4` in the PR, but double-check `@v7` is reverted to what is correct, actually the comment says "current releases include upload-artifact@v4 ... Use current supported versions").
    *   Remove mutating `go get -u ./... && go mod tidy` from monthly maintenance, change it to non-mutating freshness validation.
7.  **Preserve Snapshot behavior:**
    *   Make sure `test`, `rc`, and `alpha` triggers snapshot behavior in the GoReleaser job.
8.  **Re-verify everything:**
    *   Run tests, vet, build, and `git diff --check`.
