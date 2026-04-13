#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

INPUT="icon.png"
OUTDIR="icons"

echo " 🖼️ สร้างภาพ..."
sleep 1
[ -f "icon.png" ] || { echo "❌ icon.png missing"; exit 1; }

mkdir -p $OUTDIR

SIZES=(512 256 128 64 32 16)

for SIZE in "${SIZES[@]}"; do
  convert "$INPUT" \
    -resize ${SIZE}x${SIZE} \
    "$OUTDIR/icon-${SIZE}.png"

done
echo "✅ เสร็จแล้ว!"
