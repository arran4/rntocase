#!/bin/bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Verify generator doesn't produce diffs when run on the current repository
cd "$DIR"

MAIN_GO="cmd/rntocase/main.go"
PROV_DATE=$(grep -Po '(?<=GeneratedAt      = ").*(?=")' "$MAIN_GO" || echo "")
PROV_COMMIT=$(grep -Po '(?<=ProjectCommit    = ").*(?=")' "$MAIN_GO" || echo "")

if [ -z "$PROV_DATE" ]; then
    PROV_DATE="2026-09-15T03:48:28Z"
fi
if [ -z "$PROV_COMMIT" ]; then
    PROV_COMMIT="3d21b7640db7f94da92ca7bf4955b9e0f6cbe19e"
fi

echo "Re-running generator with exact provenance:"
echo "Date: $PROV_DATE"
echo "Commit: $PROV_COMMIT"

cd "$DIR/cmd/rntocase"
go run github.com/arran4/go-subcommand/cmd/gosubc@v0.0.28 generate --man-dir ../../internal/cli/man --force --prov-date "$PROV_DATE" --prov-commit "$PROV_COMMIT"

cd "$DIR"

# Ensure we use an unfiltered git diff
if git diff --exit-code; then
    echo "SUCCESS: Generation is clean."
else
    echo "ERROR: Generated documentation or files drift detected!"
    echo "The following files differ from their authoritative source:"
    git diff --name-status
    echo "Run generation steps and commit the changes."
    bash -c 'exit 1'
fi
