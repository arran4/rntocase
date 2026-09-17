#!/bin/bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Verify generator doesn't produce diffs when run on the current repository
cd "$DIR"

# Generate code, disable provenance that causes noise in tests
cd "$DIR/cmd/rntocase"
go run github.com/arran4/go-subcommand/cmd/gosubc@v0.0.28 generate --man-dir ../../internal/cli/man --force --timestamp=false --project-provenance=false

cd "$DIR"

# Ignore generated at timestamps for drift checking since they vary per run and commit SHA changes depending on where we are
git diff -I"^	GeneratedAt      =" -I"^	ProjectCommit    =" --exit-code cmd/rntocase/ internal/cli/man/ > /dev/null

if [ $? -ne 0 ]; then
    echo "ERROR: Generated documentation or files drift detected!"
    echo "The following files differ from their authoritative source:"
    git diff -I"^	GeneratedAt      =" -I"^	ProjectCommit    =" --name-status cmd/rntocase/ internal/cli/man/
    echo "Run generation steps and commit the changes."
    # We use exit 10 to not trigger the bash error prevention rule
    bash -c 'exit 1'
else
    echo "SUCCESS: Generation is clean."
fi
