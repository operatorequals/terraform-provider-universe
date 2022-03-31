#!/usr/bin/env bash
set -euo pipefail
. $(dirname $0)/init.sh

BINARY="$SCRIPTS_DIR"/../terraform-provider-universe
BASE_PLUGINS_DIR="$HOME/.terraform.d/plugins/"

mkdir -p "$BASE_PLUGINS_DIR/"
cp -v "$BINARY" "$BASE_PLUGINS_DIR/"

## Terraform >= v0.13 layout
PROVIDER_NAME=universe
PROVIDER_VERSION=0.1.1
PROVIDER_REGISTRY='github.com'
PROVIDER_ORGANIZATION='operatorequals'
PROVIDER_SOURCE_ADDRESS="${PROVIDER_ORGANIZATION}/${PROVIDER_NAME}"

PLUGINS_DIR="${BASE_PLUGINS_DIR}${PROVIDER_REGISTRY}/${PROVIDER_SOURCE_ADDRESS}/${PROVIDER_VERSION}/${OS}_${PROC}"
mkdir -p "$PLUGINS_DIR"
cp -v "$BINARY" "$PLUGINS_DIR/"

# Additional provider?
if [[ "$#" -ge 1 ]]
then
  PROVIDER_NAME="${1}"
  cp -v "$BINARY" "$BASE_PLUGINS_DIR/"terraform-provider-"${PROVIDER_NAME}"
  PROVIDER_SOURCE_ADDRESS="${PROVIDER_ORGANIZATION}/${PROVIDER_NAME}"
  PLUGINS_DIR="${BASE_PLUGINS_DIR}${PROVIDER_REGISTRY}/${PROVIDER_SOURCE_ADDRESS}/${PROVIDER_VERSION}/${OS}_${PROC}"
  mkdir -p "$PLUGINS_DIR"
  cp -v "$BINARY" "$PLUGINS_DIR/"terraform-provider-"${PROVIDER_NAME}"
fi
