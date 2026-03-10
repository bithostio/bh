#!/usr/bin/env bash
set -euo pipefail

REPO="bithostio/bh"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="bh"

# Detect OS.
OS="$(uname -s)"
case "$OS" in
  Darwin)  OS="darwin" ;;
  Linux)   OS="linux" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture.
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Windows gets an .exe suffix.
SUFFIX=""
if [ "$OS" = "windows" ]; then
  SUFFIX=".exe"
fi

ASSET="${BINARY_NAME}-${OS}-${ARCH}${SUFFIX}"

# Fetch latest release tag from GitHub API.
echo "Fetching latest release..."
TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4)"

if [ -z "$TAG" ]; then
  echo "Error: could not determine latest release"
  exit 1
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

echo "Downloading ${BINARY_NAME} ${TAG} for ${OS}/${ARCH}..."
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL -o "${TMPDIR}/${BINARY_NAME}${SUFFIX}" "$URL"
chmod +x "${TMPDIR}/${BINARY_NAME}${SUFFIX}"

mkdir -p "$INSTALL_DIR"
install -m 755 "${TMPDIR}/${BINARY_NAME}${SUFFIX}" "${INSTALL_DIR}/${BINARY_NAME}${SUFFIX}"

echo "Installed ${BINARY_NAME} ${TAG} to ${INSTALL_DIR}/${BINARY_NAME}${SUFFIX}"

if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo ""
  echo "Add ${INSTALL_DIR} to your PATH:"
  echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
