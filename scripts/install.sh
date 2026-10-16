#!/usr/bin/env bash
set -euo pipefail

REPO="bithostio/bh"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="bh"

# Keep redirected output and terminals without color support readable.
BOLD='' RESET='' ACCENT='' GREEN='' MUTED='' RED=''
if [ -t 1 ] && [ "${TERM:-dumb}" != "dumb" ] && [ -z "${NO_COLOR+x}" ]; then
  BOLD=$'\033[1m'
  RESET=$'\033[0m'
  ACCENT=$'\033[35m'
  GREEN=$'\033[32m'
  MUTED=$'\033[90m'
  RED=$'\033[31m'
fi

fail() {
  printf '\n  %sError:%s %s\n\n' "$RED" "$RESET" "$1" >&2
  exit 1
}

printf '\n%s%s' "$BOLD" "$ACCENT"
cat <<'BANNER'
   _     _ _   _               _
  | |__ (_) |_| |__   ___  ___| |_
  | '_ \| | __| '_ \ / _ \/ __| __|
  | |_) | | |_| | | | (_) \__ \ |_
  |_.__/|_|\__|_| |_|\___/|___/\__|
BANNER
printf '%s\n  %sBithost.io CLI%s\n' "$RESET" "$BOLD" "$RESET"
printf '  %sCreate and manage your cloud servers.%s\n\n' "$MUTED" "$RESET"

command -v curl >/dev/null 2>&1 || fail 'curl is required. Install it and run this script again.'

# Detect OS.
OS="$(uname -s)"
case "$OS" in
  Darwin)  OS="darwin" ;;
  Linux)   OS="linux" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *) fail "Unsupported OS: $OS" ;;
esac

# Detect architecture.
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *) fail "Unsupported architecture: $ARCH" ;;
esac

# Windows gets an .exe suffix.
SUFFIX=""
if [ "$OS" = "windows" ]; then
  SUFFIX=".exe"
fi

ASSET="${BINARY_NAME}-${OS}-${ARCH}${SUFFIX}"

# Fetch latest release tag from GitHub API.
printf '  %s[1/3]%s Finding the latest release...\n' "$ACCENT" "$RESET"
TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d '"' -f 4)" \
  || fail 'Could not fetch the latest release. Check your connection and try again.'

if [ -z "$TAG" ]; then
  fail 'Could not determine the latest release.'
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

printf '  %s[2/3]%s Downloading %s %s for %s/%s...\n' "$ACCENT" "$RESET" "$BINARY_NAME" "$TAG" "$OS" "$ARCH"
INSTALL_TMP_DIR="$(mktemp -d)" || fail 'Could not create a temporary directory.'
trap 'rm -rf "$INSTALL_TMP_DIR"' EXIT

curl -fsSL -o "${INSTALL_TMP_DIR}/${BINARY_NAME}${SUFFIX}" "$URL" \
  || fail "Could not download ${ASSET}. Check your connection and try again."

printf '  %s[3/3]%s Installing to %s...\n' "$ACCENT" "$RESET" "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR" || fail "Could not create ${INSTALL_DIR}. Check directory permissions."
install -m 755 "${INSTALL_TMP_DIR}/${BINARY_NAME}${SUFFIX}" "${INSTALL_DIR}/${BINARY_NAME}${SUFFIX}" \
  || fail "Could not install to ${INSTALL_DIR}. Check directory permissions."

printf '\n  %s%sInstalled bh %s successfully.%s\n' "$BOLD" "$GREEN" "$TAG" "$RESET"
printf '  %s%s/%s%s%s\n' "$MUTED" "$INSTALL_DIR" "$BINARY_NAME" "$SUFFIX" "$RESET"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    printf '\n  Add bh to your PATH before getting started:\n\n'
    printf '    %sexport PATH="%s:$PATH"%s\n' "$BOLD" "$INSTALL_DIR" "$RESET"
    printf '\n  %sSave this line in your shell profile to keep it across sessions.%s\n' "$MUTED" "$RESET"
    ;;
esac

printf '\n  %sGet started%s\n\n' "$BOLD" "$RESET"
printf '    %sbh auth%s     Connect your Bithost account\n' "$ACCENT" "$RESET"
printf '    %sbh manage%s   Open the interactive dashboard\n\n' "$ACCENT" "$RESET"
