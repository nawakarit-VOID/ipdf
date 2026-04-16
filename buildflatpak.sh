#!/bin/bash
set -e
export PATH=/usr/local/go/bin:$PATH

echo "ตรวจเช็คไฟล์"
sleep 1

[ -f "main.go" ] || { echo "❌ main.go missing"; exit 1; }
[ -f "go.mod" ] || { echo "❌ go.mod missing"; exit 1; }
[ -f "go.sum" ] || { echo "❌ go.sum, missing"; exit 1; }

echo "🔨 ใจเย็นๆ..."
sleep 1
go mod tidy
go build -ldflags="-s -w" -o ipdf

echo "📦 ลบ Flatpak repo...อันก่อน"
sleep 1
rm -rf build-dir repo

echo "📦 ก็อป iconsไปที่ Flatpak..."
cp -r icons  flatpak/

echo "📦 ก็อป assets ไปที่ Flatpak..."
cp -r assets flatpak/

echo "📦 หาไฟล์ json แล้ว รัน scrip...สร้าง ไฟล์ bundle"
sleep 1
flatpak-builder --force-clean --repo=repo build-dir flatpak/com.nawakarit.ipdf.json

echo "📦 Packing...เตรียมรับแรงกระแทก"
sleep 1
flatpak build-bundle repo ipdf.flatpak com.nawakarit.ipdf

done
echo "✅ เสร็จแล้ว!"



























