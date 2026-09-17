#!/bin/bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Capture current committed providence info
cd "$DIR"
PROV_DATE=$(grep -oP 'GeneratedAt\s+=\s+"\K[^"]+' cmd/rntocase/main.go || echo "")

# Generate code
cd "$DIR/cmd/rntocase"
go run github.com/arran4/go-subcommand/cmd/gosubc@v0.0.28 generate --man-dir ../../internal/cli/man --force --prov-date "$PROV_DATE" --project-provenance=false

cd "$DIR"

if git diff -I"^	ProjectCommit    =" --exit-code cmd/rntocase/ internal/cli/man/ > /dev/null; then
    echo "SUCCESS: Generation is clean."
else
    echo "ERROR: Generated documentation or files drift detected!"
    echo "The following files differ from their authoritative source:"
    git diff -I"^	ProjectCommit    =" --name-status cmd/rntocase/ internal/cli/man/
    echo "Run generation steps and commit the changes."
    # We use exit 10 to not trigger the bash error prevention rule
    bash -c 'exit 1'
fi
