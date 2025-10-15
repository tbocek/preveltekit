#!/bin/bash
set -e

VERSION_TYPE="patch"

show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --patch     Increment patch version (default)"
    echo "  --minor     Increment minor version"
    echo "  --major     Increment major version"
    echo "  -h, --help  Show this help message"
}

while [[ $# -gt 0 ]]; do
    case $1 in
        --patch) VERSION_TYPE="patch"; shift ;;
        --minor) VERSION_TYPE="minor"; shift ;;
        --major) VERSION_TYPE="major"; shift ;;
        -h|--help) show_usage; exit 0 ;;
        *) echo "Unknown option: $1"; show_usage; exit 1 ;;
    esac
done

increment_version() {
    local version=${1#v}
    IFS='.' read -ra PARTS <<< "$version"
    local major=${PARTS[0]}
    local minor=${PARTS[1]:-0}
    local patch=${PARTS[2]:-0}
    
    case $2 in
        major) major=$((major + 1)); minor=0; patch=0 ;;
        minor) minor=$((minor + 1)); patch=0 ;;
        patch) patch=$((patch + 1)) ;;
    esac
    
    echo "$major.$minor.$patch"
}

# Check git status
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not in a git repository"
    exit 1
fi

if [[ -n $(git status --porcelain) ]]; then
    echo "Error: Uncommitted changes. Commit first."
    git status --porcelain
    exit 1
fi

# Get current version from package.json
CURRENT_VERSION=$(node -p "require('./package.json').version")
echo "Current version: $CURRENT_VERSION"

# Calculate new version
NEW_VERSION=$(increment_version "$CURRENT_VERSION" "$VERSION_TYPE")
echo "New version: $NEW_VERSION"

read -p "Proceed with version $NEW_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# Update package.json
pnpm version "$NEW_VERSION" --no-git-tag-version
git add package.json

FILES=("README.md" "example/src/Landing.svelte" "example/src/Documentation.svelte")

# Update version in all files
for file in "${FILES[@]}"; do
    sed -i "s/\"preveltekit\": \"\\^[0-9]\+\.[0-9]\+\.[0-9]\+\"/\"preveltekit\": \"^$NEW_VERSION\"/g" "$file"
    git add "$file"
done

pnpm run docs
git add docs/**
git commit -m "release: v$NEW_VERSION"
git tag "v$NEW_VERSION"
git push && git push --tags

# Run release
pnpm run release

echo "Released v$NEW_VERSION"