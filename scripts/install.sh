#!/usr/bin/env bash
# Copyright 2025 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0


set -euo pipefail

REPO="github.com/eat-pray-ai/yutu"

ARCH="UNKNOWN"
case $(uname -m) in
  x86_64 | amd64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
esac

if [[ "${ARCH}" == "UNKNOWN" ]]; then
  echo "Unlisted architecture: $(uname -m)"
  echo "Please create an issue at ${REPO}/issues"
  exit 1
fi

OS="UNKNOWN"
case $(uname -s) in
  Linux) OS="linux" ;;
  Darwin) OS="darwin" ;;
esac

if [[ "${OS}" == "UNKNOWN" ]]; then
  echo "Unlisted OS: $(uname -s)"
  echo "Please create an issue at ${REPO}/issues"
  exit 1
fi

FILE="yutu-${OS}-${ARCH}"
curl -sSfL https://${REPO}/releases/latest/download/${FILE} -o ./yutu
chmod +x ./yutu
./yutu version

echo """
yutu🐰 is downloaded to the current directory.
You may want to move it to a directory in your PATH, e.g.:
  sudo mv ./yutu /usr/local/bin/yutu
"""
