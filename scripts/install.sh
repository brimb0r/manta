#!/usr/bin/env bash
#
# Install a terraform-provider-manta release zip for local use.
#
# Usage:
#   ./install.sh terraform-provider-manta_0.0.1_linux_amd64.zip
#
# What it does:
#   1. Extracts the binary from the zip
#   2. Places it in the Terraform filesystem_mirror plugin directory
#      so that `terraform init` picks it up locally

set -euo pipefail

HOSTNAME="registry.terraform.io"
NAMESPACE="gagno"
PROVIDER_TYPE="manta"

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <zip-file>"
  exit 1
fi

ZIP_FILE="$1"

if [[ ! -f "$ZIP_FILE" ]]; then
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
  if [[ -f "${SCRIPT_DIR}/${ZIP_FILE}" ]]; then
    ZIP_FILE="${SCRIPT_DIR}/${ZIP_FILE}"
  else
    echo "Error: zip file not found: $ZIP_FILE"
    exit 1
  fi
fi

# Parse version, os, arch from filename: terraform-provider-manta_0.0.1_linux_amd64.zip
BASENAME="$(basename "$ZIP_FILE" .zip)"
if [[ "$BASENAME" =~ ^terraform-provider-manta_([0-9]+\.[0-9]+\.[0-9]+)_([a-z]+)_([a-z0-9]+)$ ]]; then
  VERSION="${BASH_REMATCH[1]}"
  OS="${BASH_REMATCH[2]}"
  ARCH="${BASH_REMATCH[3]}"
else
  echo "Error: could not parse version/os/arch from filename: $ZIP_FILE"
  exit 1
fi

# Build the filesystem_mirror plugin directory path
case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    PLUGIN_DIR="${APPDATA}/terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PROVIDER_TYPE}/${VERSION}/${OS}_${ARCH}"
    ;;
  *)
    PLUGIN_DIR="${HOME}/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${PROVIDER_TYPE}/${VERSION}/${OS}_${ARCH}"
    ;;
esac

mkdir -p "$PLUGIN_DIR"

# Extract to a temp directory first, then copy the binary into the plugin dir
TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "Extracting ${ZIP_FILE} ..."
unzip -o "$ZIP_FILE" -d "$TEMP_DIR"

# Copy the binary into the plugin directory
BINARY="$(find "$TEMP_DIR" -name 'terraform-provider-manta*' -type f | head -1)"
if [[ -z "$BINARY" ]]; then
  echo "Error: no terraform-provider-manta binary found in zip"
  exit 1
fi

cp "$BINARY" "$PLUGIN_DIR/"
chmod +x "$PLUGIN_DIR/$(basename "$BINARY")"

echo "Installed to: ${PLUGIN_DIR}/$(basename "$BINARY")"
echo ""
echo "Done! v${VERSION} is ready for local use."
echo "Run 'terraform init' to pick up the local provider."
