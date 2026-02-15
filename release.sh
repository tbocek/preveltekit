#!/bin/bash

set -e

# Default values
VERSION_TYPE="patch"

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --patch     Increment patch version (default)"
    echo "  --minor     Increment minor version"
    echo "  --major     Increment major version"
    echo "  -h, --help  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0              # Increment patch version"
    echo "  $0 --minor      # Increment minor version"
    echo "  $0 --major      # Increment major version"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --patch)
            VERSION_TYPE="patch"
            shift
            ;;
        --minor)
            VERSION_TYPE="minor"
            shift
            ;;
        --major)
            VERSION_TYPE="major"
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Function to increment version
increment_version() {
    local version=$1
    local type=$2

    # Remove 'v' prefix if present
    version=${version#v}

    # Split version into parts
    IFS='.' read -ra PARTS <<< "$version"
    major=${PARTS[0]}
    minor=${PARTS[1]:-0}
    patch=${PARTS[2]:-0}

    case $type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
    esac

    echo "$major.$minor.$patch"
}

# Build site and copy to docs/
echo "Building site..."
(cd site && bash build.sh --release)
echo "Copying site/dist to docs/..."
rm -rf docs
cp -r site/dist docs

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not in a git repository"
    exit 1
fi

# Check for uncommitted changes (including untracked files)
if [[ -n $(git status --porcelain) ]]; then
    echo "Error: You have uncommitted changes or untracked files. Please commit all changes before creating a release."
    echo ""
    echo "Uncommitted/untracked files:"
    git status --porcelain
    exit 1
fi

# Check for unpushed commits
LOCAL_COMMIT=$(git rev-parse HEAD)
REMOTE_COMMIT=$(git rev-parse @{u} 2>/dev/null || echo "")

if [[ -n "$REMOTE_COMMIT" && "$LOCAL_COMMIT" != "$REMOTE_COMMIT" ]]; then
    # Check if we're ahead of remote
    if git merge-base --is-ancestor "$REMOTE_COMMIT" "$LOCAL_COMMIT"; then
        echo "Error: You have unpushed commits. Please push all commits before creating a release."
        echo ""
        echo "Unpushed commits:"
        git log --oneline "$REMOTE_COMMIT..HEAD"
        exit 1
    else
        echo "Warning: Your local branch is behind remote. Consider pulling latest changes."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Cancelled."
            exit 0
        fi
    fi
fi

# Fetch all tags from remote
echo "Fetching tags from remote..."
git fetch --tags

# Get the latest tag
echo "Finding latest tag..."
LATEST_TAG=$(git tag --sort=-version:refname | head -n1)

if [[ -z "$LATEST_TAG" ]]; then
    echo "No existing tags found. Starting from v1.0.0"
    LATEST_TAG="v1.0.0"
    NEW_TAG="v1.0.0"
else
    echo "Latest tag: $LATEST_TAG"
    NEW_VERSION=$(increment_version "$LATEST_TAG" "$VERSION_TYPE")
    NEW_TAG="v$NEW_VERSION"
fi

echo "New version will be: $NEW_TAG (incrementing $VERSION_TYPE)"

# Confirm before proceeding
read -p "Create and push tag $NEW_TAG? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

# Create the new tag
echo "Creating tag $NEW_TAG..."
git tag -a "$NEW_TAG" -m "Release $NEW_TAG"

# Push the tag to remote
echo "Pushing tag to remote..."
git push origin "$NEW_TAG"

echo "Successfully created and pushed release $NEW_TAG"
