#!/bin/bash
# build.sh — Build RuntimeGuard binaries for Linux and Windows in the shared folder

SHARED_FOLDER="/mnt/shared/runxguard/dist"
mkdir -p "$SHARED_FOLDER"

echo "Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -o "$SHARED_FOLDER/runtimeguard" main.go
if [ $? -eq 0 ]; then
    echo "Linux binary built successfully: $SHARED_FOLDER/runtimeguard"
else
    echo "Failed to build Linux binary."
fi

echo "Building Windows binary..."
GOOS=windows GOARCH=amd64 go build -o "$SHARED_FOLDER/runtimeguard.exe" main.go
if [ $? -eq 0 ]; then
    echo "Windows binary built successfully: $SHARED_FOLDER/runtimeguard.exe"
else
    echo "Failed to build Windows binary."
fi

echo "Binaries in shared folder:"
ls -lh "$SHARED_FOLDER"

