#!/usr/bin/env bash
#
# build-fluent-zip.sh
# Resizes all 3D Fluent emoji PNGs to 256x256 and packages them into fluent.zip.
#

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Ensure the fluent repo is available
source "$ROOT/scripts/fetch-fluent-repo.sh"

OUTDIR="$ROOT/tmp/fluent"
ZIPFILE="$ROOT/fluent.zip"

rm -rf "$OUTDIR"
mkdir -p "$OUTDIR"

echo "Resizing 3D PNGs to 128x128..."

find "$FLUENT_REPO_DIR/assets" -type f -name '*_3d.png' -print0 | while IFS= read -r -d '' src; do
  name="$(basename "$src")"
  convert "$src" -resize 128x128 "$OUTDIR/$name"
done

COUNT=$(find "$OUTDIR" -type f -name '*.png' | wc -l)
echo "Resized $COUNT PNGs."

echo "Creating fluent.zip..."
(cd "$OUTDIR" && zip -q "$ZIPFILE" *.png)

echo "Done: $ZIPFILE ($(du -h "$ZIPFILE" | cut -f1))"
