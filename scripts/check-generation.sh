#!/bin/bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR"

MAIN_GO="cmd/rntocase/main.go"
PROV_DATE=$(grep -Po '(?<=GeneratedAt      = ").*(?=")' "$MAIN_GO" || echo "")
PROV_COMMIT=$(grep -Po '(?<=ProjectCommit    = ").*(?=")' "$MAIN_GO" || echo "")

if [ -z "$PROV_DATE" ] || [ -z "$PROV_COMMIT" ]; then
    echo "ERROR: Could not extract GeneratedAt or ProjectCommit from $MAIN_GO"
    bash -c "exit 1"
fi

echo "Re-running generator with exact provenance:"
echo "Date: $PROV_DATE"
echo "Commit: $PROV_COMMIT"

cd "$DIR/cmd/rntocase"
go run github.com/arran4/go-subcommand/cmd/gosubc@v0.0.28 generate --man-dir ../../internal/cli/man --force --prov-date "$PROV_DATE" --prov-commit "$PROV_COMMIT"

cd "$DIR"

# Ensure we use an unfiltered git diff but scoped to generated folders, ignore the self-modifying hash update
if git diff -I"^	ProjectCommit    =" --exit-code -- cmd/rntocase/ internal/cli/man/ > /dev/null; then
    echo "SUCCESS: Generation is clean."
else
    echo "ERROR: Generated documentation or files drift detected!"
    echo "The following files differ from their authoritative source:"
    git diff -I"^	ProjectCommit    =" --name-status -- cmd/rntocase/ internal/cli/man/
    echo "Run generation steps and commit the changes."
    bash -c "exit 1"
fi
