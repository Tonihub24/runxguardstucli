#!/bin/bash
# =========================================
# RuntimeGuard Full Setup Script for Linux
# =========================================
# Author: Antonio Kione
# Purpose: One-shot installation and setup for RuntimeGuard CLI
# =========================================

set -e  # Exit on any error
echo "🔹 Starting RuntimeGuard setup..."

# 1️⃣ Install Go 1.26.1
GO_VERSION="1.26.1"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"

echo "➡️ Installing Go ${GO_VERSION}..."
wget -q $GO_URL -O /tmp/$GO_TAR
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/$GO_TAR
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin
go version

# 2️⃣ Clone RuntimeGuard repository
REPO_DIR="$HOME/scripts/projects/runxguardstucli"
if [ ! -d "$REPO_DIR" ]; then
    echo "➡️ Cloning RuntimeGuard repository..."
    mkdir -p "$(dirname $REPO_DIR)"
    git clone https://github.com/<your-user>/runxguardstucli.git $REPO_DIR
else
    echo "➡️ RuntimeGuard directory exists, pulling latest updates..."
    cd $REPO_DIR
    git pull
fi
cd $REPO_DIR

# 3️⃣ Build binaries
echo "➡️ Building Linux & Windows binaries..."
chmod +x build.sh
./build.sh

# 4️⃣ Setup baseline file
echo "➡️ Setting up baseline file..."
mkdir -p ~/.runtimeguard
if [ -f "$REPO_DIR/runtimeguard_baseline.json" ]; then
    cp "$REPO_DIR/runtimeguard_baseline.json" ~/.runtimeguard/baseline.json
    echo "✅ Baseline file copied to ~/.runtimeguard/baseline.json"
else
    echo "⚠️ No baseline file found, you can generate one with:"
    echo "   ./runtimeguard init"
fi

# 5️⃣ Optional: make Linux binary global
echo "➡️ Installing Linux CLI globally..."
sudo mv "$REPO_DIR/build/runtimeguard" /usr/local/bin/
chmod +x /usr/local/bin/runtimeguard

# 6️⃣ Test installation
echo "➡️ Testing RuntimeGuard CLI..."
runtimeguard help

echo "🎉 RuntimeGuard setup completed successfully!"
