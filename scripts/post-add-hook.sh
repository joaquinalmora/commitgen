#!/bin/sh
# commitgen post-commit hook
# Automatically cache commit message for staged changes

# Only run if commitgen is available
if ! command -v commitgen >/dev/null 2>&1; then
    # Try to find commitgen in common locations
    if [ -f "./bin/commitgen" ]; then
        COMMITGEN="./bin/commitgen"
    elif [ -f "$(git rev-parse --show-toplevel)/bin/commitgen" ]; then
        COMMITGEN="$(git rev-parse --show-toplevel)/bin/commitgen"
    else
        exit 0
    fi
else
    COMMITGEN="commitgen"
fi

# Check if there are staged changes
if git diff --cached --quiet; then
    exit 0
fi

# Generate cache in background (don't slow down git add)
$COMMITGEN cache --verbose >/dev/null 2>&1 &
