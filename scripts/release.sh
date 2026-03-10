#!/usr/bin/env bash
set -euo pipefail

REPO="bithostio/bh"

usage() {
  echo "Usage: $0 <version>"
  echo "  version: semver tag, e.g. v0.1.0"
  exit 1
}

if [ $# -ne 1 ]; then
  usage
fi

VERSION="$1"

# Validate version format.
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: version must match v*.*.*  (e.g. v0.1.0)"
  exit 1
fi

echo "==> Running tests..."
make test

echo "==> Tagging $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

echo "==> Building all platforms..."
make clean build-all VERSION="$VERSION"

echo "==> Generating checksums..."
(cd dist && shasum -a 256 bh-* > bh-checksums.txt)

echo "==> Creating GitHub Release..."
gh release create "$VERSION" \
  --repo "$REPO" \
  --title "$VERSION" \
  --generate-notes \
  dist/*

echo "==> Released $VERSION"
echo "    https://github.com/$REPO/releases/tag/$VERSION"
