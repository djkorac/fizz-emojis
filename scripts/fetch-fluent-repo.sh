#!/usr/bin/env bash
#
# fetch-fluent-repo.sh
# Handles the shallow clone and sparse checkout of the Microsoft Fluent UI Emoji repo.
#

set -euo pipefail

# Define paths relative to the script location
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
REPO_URL="https://github.com/microsoft/fluentui-emoji.git"
REPO_DIR="$ROOT/tmp/fluent-repo"

echo "Ensuring Fluent UI Emoji repository is available..."

if [ -d "$REPO_DIR/.git" ]; then
  echo "Repository already exists at $REPO_DIR. skipping clone."
else
  echo "Cloning repository (shallow, sparse)..."
  mkdir -p "$REPO_DIR"

  # Clone with blobless filter for speed
  git clone --depth=1 --filter=blob:none --sparse "$REPO_URL" "$REPO_DIR"

  pushd "$REPO_DIR" > /dev/null
  # Set sparse checkout to only pull the assets folder
  git sparse-checkout set assets
  popd > /dev/null
fi

# Export the directory path so other scripts can use it if they source this
export FLUENT_REPO_DIR="$REPO_DIR"

echo "Repo ready at: $FLUENT_REPO_DIR"
