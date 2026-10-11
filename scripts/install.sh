#!/usr/bin/env bash
# Copyright 2025 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0


set -euo pipefail

REPO="github.com/eat-pray-ai/yutu"
INSTALL_DIR="${XDG_BIN_HOME:-${HOME}/.local/bin}"

ARCH="UNKNOWN"
case $(uname -m) in
  x86_64 | amd64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
esac

if [[ "${ARCH}" == "UNKNOWN" ]]; then
  echo "Unsupported architecture: $(uname -m)"
  echo "Report: https://${REPO}/issues"
  exit 1
fi

OS="UNKNOWN"
case $(uname -s) in
  Linux) OS="linux" ;;
  Darwin) OS="darwin" ;;
esac

if [[ "${OS}" == "UNKNOWN" ]]; then
  echo "Unsupported OS: $(uname -s)"
  echo "Report: https://${REPO}/issues"
  exit 1
fi

FILE="yutu-${OS}-${ARCH}"
mkdir -p "${INSTALL_DIR}"
curl -sSfL https://${REPO}/releases/latest/download/${FILE} -o "${INSTALL_DIR}/yutu"
chmod +x "${INSTALL_DIR}/yutu"
"${INSTALL_DIR}/yutu" version

if command -v npx >/dev/null 2>&1; then
  echo "Installing skill..."
  if npx --yes skills add eat-pray-ai/yutu --skill youtube -y; then
    echo "Skill installed."
  else
    echo "Skill install failed. Run manually:"
    echo "  npx skills add eat-pray-ai/yutu --skill youtube"
  fi
else
  echo "Skill skipped: npx not found."
fi

echo """
yutu🐰 installed:
  ${INSTALL_DIR}/yutu

Add ${INSTALL_DIR} to PATH if needed.
"""
