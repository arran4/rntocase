import sys
import re

with open('.github/workflows/ci.yaml', 'r') as f:
    content = f.read()

# 1. Update pull_request triggers (remove reopened, ready_for_review, closed)
content = re.sub(
    r'  pull_request:\n    types: \[opened, synchronize, reopened, ready_for_review, closed\]',
    r'  pull_request:\n    types: [opened, synchronize]',
    content
)

# Set workflow_dispatch mode to required: true, default: build
content = re.sub(
    r'        type: choice\n        options:\n          - lint-fix\n          - publish-tag',
    r'        type: choice\n        required: true\n        default: build\n        options:\n          - build\n          - lint-fix\n          - publish-tag',
    content
)

# Refactor route job script
route_script_old = r'''          run_code_checks=true
          run_build=true
          run_release=false
          run_autofix=false
          run_publisher=false
          run_maintenance=false
          mode="${INPUT_MODE:-}"

          if [[ "$GITHUB_EVENT_NAME" == "schedule" ]]; then
             run_maintenance=true
             run_code_checks=false
             run_build=false
          elif [[ "$GITHUB_EVENT_NAME" == "workflow_dispatch" ]]; then
             if [[ "$mode" == "lint-fix" ]]; then
                run_autofix=true
             elif [[ "$mode" == "publish-tag" ]]; then
                run_code_checks=false
                run_build=false
                run_publisher=true
             elif [[ "$mode" == "release-"* ]]; then
                run_release=true
             fi
          elif [[ "$GITHUB_EVENT_NAME" == "push" && "$GITHUB_REF_TYPE" == "tag" && "$GITHUB_REF" == refs/tags/v* ]]; then
             run_publisher=true
          fi'''

route_script_new = r'''          run_code_checks=false
          run_build=false
          run_release=false
          run_autofix=false
          run_publisher=false
          run_maintenance=false
          mode="${INPUT_MODE:-}"

          if [[ "$GITHUB_EVENT_NAME" == "push" ]]; then
             if [[ "$GITHUB_REF_TYPE" == "tag" && "$GITHUB_REF" == refs/tags/v* ]]; then
                run_publisher=true
             else
                run_code_checks=true
                run_build=true
             fi
          elif [[ "$GITHUB_EVENT_NAME" == "pull_request" ]]; then
             run_code_checks=true
             run_build=true
          elif [[ "$GITHUB_EVENT_NAME" == "schedule" ]]; then
             run_maintenance=true
          elif [[ "$GITHUB_EVENT_NAME" == "workflow_dispatch" ]]; then
             if [[ "$mode" == "build" ]]; then
                run_code_checks=true
                run_build=true
             elif [[ "$mode" == "lint-fix" ]]; then
                run_autofix=true
             elif [[ "$mode" == "publish-tag" ]]; then
                run_code_checks=true
                run_build=true
                run_publisher=true
             elif [[ "$mode" == "release-"* ]]; then
                run_code_checks=true
                run_build=true
                run_release=true
             fi
          fi'''

content = content.replace(route_script_old, route_script_new)

# Monthly Maintenance - Make non-mutating
maintenance_old = r'''      - run: go get -u ./... && go mod tidy
      - name: Create Pull Request
        if: ${{ github.event_name == 'schedule' || inputs.allow_prs != false }}
        uses: peter-evans/create-pull-request@v8
        with:
          commit-message: "chore: monthly dependency update"
          title: "chore: monthly dependency update"
          branch: automation/maintenance
          delete-branch: true'''

maintenance_new = r'''      - name: Validate dependencies are up to date
        run: |
          go get -u ./...
          go mod tidy
          if ! git diff --quiet; then
            echo "Dependencies are outdated. Run 'go get -u ./... && go mod tidy' locally."
            sh -c "exit 1"
          fi'''
content = content.replace(maintenance_old, maintenance_new)

# Remove cleanup-autofix-prs job
cleanup_job_pattern = r'  cleanup-autofix-prs:\n    name: Cleanup Autofix PRs.*?(?=\n  build:)'
content = re.sub(cleanup_job_pattern, '', content, flags=re.DOTALL)


# release-ready gate
release_ready_old = r'''  release-ready:
    name: Release Quality Gates Passed
    needs: [route, validation, build]
    if: always() && !contains(needs.*.result, 'failure') && !contains(needs.*.result, 'cancelled')
    runs-on: ubuntu-latest
    steps:
      - run: echo "All release quality gates passed."'''

release_ready_new = r'''  release-ready:
    name: Release Quality Gates Passed
    needs: [route, validation, build]
    if: ${{ (needs.route.outputs.run_release == 'true' || needs.route.outputs.run_publisher == 'true') && needs.validation.result == 'success' && needs.build.result == 'success' }}
    runs-on: ubuntu-latest
    steps:
      - run: echo "All release quality gates passed."'''
content = content.replace(release_ready_old, release_ready_new)

# prepare-release-tag env overrides and auth branch checking
prepare_env_old = r'''        env:
          GITHUB_REF_NAME: ${{ github.ref }}
          GITHUB_SHA: ${{ github.sha }}
        run: |
          set -euo pipefail
          if [[ "$GITHUB_REF_NAME" != "refs/heads/main" ]]; then
            echo "Error: Manual release preparation must run on refs/heads/main, got $GITHUB_REF_NAME"
            sh -c "exit 1"
          fi
          git fetch origin main
          MAIN_SHA=$(git rev-parse origin/main)
          if [[ "$MAIN_SHA" != "$GITHUB_SHA" ]]; then
            echo "Error: Requested release against $GITHUB_SHA but origin/main is at $MAIN_SHA"
            sh -c "exit 1"
          fi'''

prepare_env_new = r'''        run: |
          set -euo pipefail
          AUTHORITATIVE_BRANCH="${GITHUB_REF#refs/heads/}"
          if [[ "$AUTHORITATIVE_BRANCH" != "main" ]]; then
            echo "Error: Manual release preparation must run on refs/heads/main, got $AUTHORITATIVE_BRANCH"
            sh -c "exit 1"
          fi
          git fetch origin "$AUTHORITATIVE_BRANCH"
          MAIN_SHA=$(git rev-parse "origin/$AUTHORITATIVE_BRANCH")
          if [[ "$MAIN_SHA" != "$GITHUB_SHA" ]]; then
            echo "Error: Requested release against $GITHUB_SHA but origin/$AUTHORITATIVE_BRANCH is at $MAIN_SHA"
            sh -c "exit 1"
          fi'''
content = content.replace(prepare_env_old, prepare_env_new)

# fix regex
content = content.replace(r'''^v[0-9]+.[0-9]+.[0-9]+(-[a-zA-Z0-9.]+)?$''', r'''^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$''')

# prepare-release-tag concurrency
prepare_job_old = r'''  prepare-release-tag:
    name: Prepare Release Tag
    needs: [route, release-ready]
    if: ${{ needs.route.outputs.run_release == 'true' }}
    runs-on: ubuntu-latest
    permissions:
      contents: write
      actions: write'''
prepare_job_new = r'''  prepare-release-tag:
    name: Prepare Release Tag
    needs: [route, release-ready]
    if: ${{ needs.route.outputs.run_release == 'true' }}
    runs-on: ubuntu-latest
    concurrency:
      group: ${{ github.workflow }}-release-preparation
      cancel-in-progress: false
    permissions:
      contents: write
      actions: write'''
content = content.replace(prepare_job_old, prepare_job_new)


# publisher concurrency and goreleaser args
publisher_job_old = r'''  publisher:
    name: Release Publisher
    needs: [route, release-ready]
    if: ${{ needs.route.outputs.run_publisher == 'true' }}
    runs-on: ubuntu-latest
    permissions:
      contents: write
      packages: write'''

publisher_job_new = r'''  publisher:
    name: Release Publisher
    needs: [route, release-ready]
    if: ${{ needs.route.outputs.run_publisher == 'true' }}
    runs-on: ubuntu-latest
    concurrency:
      group: ${{ github.workflow }}-publish-${{ github.ref }}
      cancel-in-progress: false
    permissions:
      contents: write
      packages: write'''
content = content.replace(publisher_job_old, publisher_job_new)


# GoReleaser arguments snapshot handling and version
goreleaser_args_old = r'''          if [[ "$GITHUB_REF" == refs/tags/*-test* || "$GITHUB_REF" == refs/tags/test-* ]]; then
             echo "args=release --snapshot --clean" >> "$GITHUB_OUTPUT"
          else
             echo "args=release --clean" >> "$GITHUB_OUTPUT"
          fi
      - name: Run GoReleaser (Sole Release Owner)
        uses: goreleaser/goreleaser-action@v7
        with:
          distribution: goreleaser
          version: latest
          args: ${{ steps.args.outputs.args }}'''

goreleaser_args_new = r'''          if [[ "$GITHUB_REF" == refs/tags/*-test* || "$GITHUB_REF" == refs/tags/test-* || "$GITHUB_REF" == refs/tags/*-rc* || "$GITHUB_REF" == refs/tags/*-alpha* ]]; then
             echo "args=release --snapshot --clean" >> "$GITHUB_OUTPUT"
          else
             echo "args=release --clean" >> "$GITHUB_OUTPUT"
          fi
      - name: Run GoReleaser (Sole Release Owner)
        uses: goreleaser/goreleaser-action@v7
        with:
          distribution: goreleaser
          version: '~> v2'
          args: ${{ steps.args.outputs.args }}'''
content = content.replace(goreleaser_args_old, goreleaser_args_new)


with open('.github/workflows/ci.yaml', 'w') as f:
    f.write(content)
