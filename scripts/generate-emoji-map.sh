#!/usr/bin/env bash
#
# generate-emoji-map.sh
# Uses the fetched repo to extract 3D emoji names and generate the JSON map.
#

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUTFILE="$ROOT/emojis.json"

# Fetch the repo (or verify it exists)
source "$ROOT/scripts/fetch-fluent-repo.sh"

echo "Extracting 3D emoji names and generating map..."

# Run the processing logic using the variable exported by the fetch script
# We find all files ending in _3d.png
find "$FLUENT_REPO_DIR/assets" -type f -name '*_3d.png' -printf '%f\n' \
  | sed 's/_3d\.png$//' \
  | sort \
  | go run "$ROOT/scripts/generate-emoji-map/main.go" -o "$OUTFILE"

echo "Done: $OUTFILE"
