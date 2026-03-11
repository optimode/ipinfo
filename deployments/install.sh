#!/usr/bin/env bash
set -euo pipefail

REPO="optimode/ipinfo"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *) echo "Error: unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *) echo "Error: unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Determine version (latest release or user-specified)
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
  if [ -z "$VERSION" ]; then
    echo "Error: could not determine latest version" >&2
    exit 1
  fi
fi

BINARY="ipinfo-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ipinfo ${VERSION} (${OS}/${ARCH})..."
curl -fsSL -o "${TMPDIR}/${BINARY}" "$DOWNLOAD_URL"

# Install binary
echo "Installing..."
chmod +x "${TMPDIR}/${BINARY}"
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/ipinfo"
  SUDO=""
else
  sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/ipinfo"
  SUDO="sudo"
fi

# Create config directory
${SUDO:-sudo} mkdir -p /etc/ipinfo

# Download example config
CONFIG_URL="https://raw.githubusercontent.com/${REPO}/${VERSION}/configs/config.example.yaml"
curl -fsSL -o "${TMPDIR}/config.example.yaml" "$CONFIG_URL"

if [ ! -f /etc/ipinfo/config.yaml ]; then
  # Fresh install: use as default config
  ${SUDO:-sudo} cp "${TMPDIR}/config.example.yaml" /etc/ipinfo/config.yaml
  echo "Configuration: /etc/ipinfo/config.yaml (new)"
else
  # Upgrade: keep existing config, save example for reference
  ${SUDO:-sudo} cp "${TMPDIR}/config.example.yaml" /etc/ipinfo/config.example.yaml
  echo "Configuration: /etc/ipinfo/config.yaml (kept existing, example saved)"
fi

echo "Installed ipinfo ${VERSION} to ${INSTALL_DIR}/ipinfo"
echo ""
echo "Next: edit /etc/ipinfo/config.yaml and set your ip-api.com Pro API key"
