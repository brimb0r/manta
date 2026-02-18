#!/usr/bin/env bash
#
# Install a terraform-provider-manta release zip for local dev use.
#
# Usage:
#   ./install.sh terraform-provider-manta_0.0.1_linux_amd64.zip
#   ./install.sh terraform-provider-manta_0.0.1_windows_amd64.zip
#
# What it does:
#   1. Extracts the binary from the zip
#   2. Puts it in a dev_overrides directory
#   3. Ensures your Terraform CLI config has a dev_overrides block
#      pointing to that directory so you can skip `terraform init`

set -euo pipefail

PROVIDER_NAMESPACE="gagno"
PROVIDER_TYPE="manta"
PROVIDER_SOURCE="registry.terraform.io/${PROVIDER_NAMESPACE}/${PROVIDER_TYPE}"

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <zip-file>"
  exit 1
fi

ZIP_FILE="$1"

if [[ ! -f "$ZIP_FILE" ]]; then
  # Try looking in the same directory as the script
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
  if [[ -f "${SCRIPT_DIR}/${ZIP_FILE}" ]]; then
    ZIP_FILE="${SCRIPT_DIR}/${ZIP_FILE}"
  else
    echo "Error: zip file not found: $ZIP_FILE"
    exit 1
  fi
fi

# Parse version from zip name: terraform-provider-manta_0.0.1_linux_amd64.zip
BASENAME="$(basename "$ZIP_FILE" .zip)"
VERSION="$(echo "$BASENAME" | sed -E 's/^terraform-provider-manta_([0-9]+\.[0-9]+\.[0-9]+).*/\1/')"
if [[ -z "$VERSION" ]]; then
  echo "Error: could not parse version from filename: $ZIP_FILE"
  exit 1
fi

# Detect OS
case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    IS_WINDOWS=true
    DEV_DIR="${APPDATA}/terraform-provider-manta-dev"
    RC_FILE="${APPDATA}/terraform.rc"
    ;;
  *)
    IS_WINDOWS=false
    DEV_DIR="${HOME}/.local/share/terraform-provider-manta-dev"
    RC_FILE="${HOME}/.terraformrc"
    ;;
esac

mkdir -p "$DEV_DIR"

# Extract
echo "Extracting ${ZIP_FILE} ..."
unzip -o "$ZIP_FILE" -d "$DEV_DIR"

# On Linux/macOS make sure the binary is executable
if [[ "$IS_WINDOWS" == false ]]; then
  chmod +x "${DEV_DIR}"/terraform-provider-manta*
fi

echo "Binary installed to: ${DEV_DIR}"

# Configure dev_overrides in the Terraform CLI config.
# dev_overrides lets Terraform use the local binary directly — no init needed.
DEV_DIR_ESCAPED="$DEV_DIR"
if [[ "$IS_WINDOWS" == true ]]; then
  # Convert backslashes to forward slashes for HCL
  DEV_DIR_ESCAPED="$(cygpath -m "$DEV_DIR" 2>/dev/null || echo "$DEV_DIR" | sed 's|\\|/|g')"
fi

if [[ -f "$RC_FILE" ]] && grep -q "$PROVIDER_SOURCE" "$RC_FILE" 2>/dev/null; then
  echo "dev_overrides already configured in ${RC_FILE}"
else
  # Append a provider_installation block with dev_overrides
  if [[ -f "$RC_FILE" ]] && grep -q "dev_overrides" "$RC_FILE" 2>/dev/null; then
    echo ""
    echo "WARNING: ${RC_FILE} already has a dev_overrides block but does not include"
    echo "  ${PROVIDER_SOURCE}"
    echo ""
    echo "Add this line inside the dev_overrides block manually:"
    echo "    \"${PROVIDER_SOURCE}\" = \"${DEV_DIR_ESCAPED}\""
  else
    cat >> "$RC_FILE" <<EOF

provider_installation {
  dev_overrides {
    "${PROVIDER_SOURCE}" = "${DEV_DIR_ESCAPED}"
  }
  direct {}
}
EOF
    echo "Wrote dev_overrides to ${RC_FILE}"
  fi
fi

echo ""
echo "Done! v${VERSION} is ready for local use."
echo "Terraform will use the local binary — no 'terraform init' required."
