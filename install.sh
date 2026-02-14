#!/bin/sh
set -e

REPO="daedaluz/mantra-cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="mantra-cli"

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)      echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)       echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

ASSET="${BINARY_NAME}-${OS}-${ARCH}"

echo "Fetching latest release from ${REPO}..."
TAG=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | head -1 | cut -d'"' -f4)

if [ -z "$TAG" ]; then
    echo "Error: could not determine latest release." >&2
    exit 1
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"
echo "Downloading ${ASSET} (${TAG})..."

TMPFILE=$(mktemp)
trap 'rm -f "$TMPFILE"' EXIT

HTTP_CODE=$(curl -sL -o "$TMPFILE" -w "%{http_code}" "$URL")
if [ "$HTTP_CODE" != "200" ]; then
    echo "Error: download failed (HTTP ${HTTP_CODE})." >&2
    echo "URL: ${URL}" >&2
    exit 1
fi

chmod +x "$TMPFILE"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPFILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Installing to ${INSTALL_DIR} (requires sudo)..."
    sudo mv "$TMPFILE" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "Installed ${BINARY_NAME} ${TAG} to ${INSTALL_DIR}/${BINARY_NAME}"
