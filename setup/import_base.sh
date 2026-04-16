#!/bin/bash

# 1. Define the ~/.docksmith storage directories
DOCKSMITH_DIR="$HOME/.docksmith"
LAYERS_DIR="$DOCKSMITH_DIR/layers"
IMAGES_DIR="$DOCKSMITH_DIR/images"

# Create the directories if they don't exist yet
mkdir -p "$LAYERS_DIR"
mkdir -p "$IMAGES_DIR"

echo "Downloading Alpine Linux base image (ARM64)..."
# URL for the official Alpine Linux minimal root filesystem
ALPINE_URL="https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/aarch64/alpine-minirootfs-3.19.1-aarch64.tar.gz"
TEMP_TAR_GZ="/tmp/alpine-rootfs.tar.gz"
TEMP_TAR="/tmp/alpine-rootfs.tar"

# Download the file using curl
curl -s -L -o "$TEMP_TAR_GZ" "$ALPINE_URL"

echo "Decompressing..."
# Unzip it into a standard .tar file to make it easier for Docksmith to process
gunzip -f "$TEMP_TAR_GZ"

echo "Calculating SHA-256 layer digest..."
# Calculate the hash (handles both Mac and Linux environments)
if command -v sha256sum >/dev/null 2>&1; then
    HASH=$(sha256sum "$TEMP_TAR" | awk '{print $1}')
else
    HASH=$(shasum -a 256 "$TEMP_TAR" | awk '{print $1}')
fi

echo "Layer Digest: $HASH"

# Move the layer into the Docksmith storage folder, named by its hash
mv "$TEMP_TAR" "$LAYERS_DIR/${HASH}.tar"

echo "Creating image manifest..."
# Create the JSON manifest file so Docksmith knows 'alpine' exists
cat <<EOF > "$IMAGES_DIR/alpine.json"
{
  "name": "alpine",
  "tag": "latest",
  "created": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "layers": [
    "$HASH"
  ],
  "env": ["PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"],
  "cmd": ["/bin/sh"]
}
EOF

echo "Success! The base image 'alpine:latest' is securely stored and ready for offline builds."
