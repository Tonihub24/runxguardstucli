#!/bin/bash
# build.sh — Build RuntimeGuard binaries for Linux and Windows in the home folder

# Use a local build folder instead of /mnt/shared
BUILD_FOLDER="$HOME/scripts/projects/runxguardstucli/build"
mkdir -p "$BUILD_FOLDER"

echo "Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -o "$BUILD_FOLDER/runtimeguard" main.go
if [ $? -eq 0 ]; then
    echo "Linux binary built successfully: $BUILD_FOLDER/runtimeguard"
else
    echo "Failed to build Linux binary."
fi

echo "Building Windows binary..."
GOOS=windows GOARCH=amd64 go build -o "$BUILD_FOLDER/runtimeguard.exe" main.go
if [ $? -eq 0 ]; then
    echo "Windows binary built successfully: $BUILD_FOLDER/runtimeguard.exe"
else
    echo "Failed to build Windows binary."
fi

echo "Binaries in build folder:"
ls -lh "$BUILD_FOLDER"
