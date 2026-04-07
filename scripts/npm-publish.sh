#!/usr/bin/env bash
# Copyright 2026 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# Usage: scripts/npm-publish.sh <version> <dist-dir> [--dry-run]
# The version should NOT have a "v" prefix.

VERSION="${1:?Usage: npm-publish.sh <version> <dist-dir> [--dry-run]}"
DIST_DIR="${2:?Usage: npm-publish.sh <version> <dist-dir> [--dry-run]}"
DRY_RUN="${3:-}"

REPO_ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
NPM_DIR="$REPO_ROOT/npm"

NPM_TAG="latest"
[[ "$VERSION" == *-* ]] && NPM_TAG="next"

PUBLISH_FLAGS=(--access public --tag "$NPM_TAG")
[[ -n "${GITHUB_ACTIONS:-}" ]] && PUBLISH_FLAGS+=(--provenance)
[[ "$DRY_RUN" == "--dry-run" ]] && PUBLISH_FLAGS+=(--dry-run)

declare -A PLATFORM_MAP=(
  ["yutu-darwin-arm64"]="yutu_darwin_arm64_v8.0/yutu-darwin-arm64"
  ["yutu-darwin-x64"]="yutu_darwin_amd64_v1/yutu-darwin-amd64"
  ["yutu-linux-arm64"]="yutu_linux_arm64_v8.0/yutu-linux-arm64"
  ["yutu-linux-x64"]="yutu_linux_amd64_v1/yutu-linux-amd64"
  ["yutu-win32-arm64"]="yutu_windows_arm64_v8.0/yutu-windows-arm64.exe"
  ["yutu-win32-x64"]="yutu_windows_amd64_v1/yutu-windows-amd64.exe"
)

stamp_version() { sed -i.bak "s/\"0.0.0\"/\"${VERSION}\"/g" "$1" && rm -f "$1.bak"; }
restore_version() { sed -i.bak "s/\"${VERSION}\"/\"0.0.0\"/g" "$1" && rm -f "$1.bak"; }

publish_platform() {
  local pkg_dir="$1" src_path="$DIST_DIR/${PLATFORM_MAP[$1]}"
  local target_dir="$NPM_DIR/$pkg_dir"
  local target_bin="yutu"
  [[ "$pkg_dir" == *win32* ]] && target_bin="yutu.exe"

  [[ -f "$src_path" ]] || { echo "ERROR: Binary not found: $src_path"; exit 1; }

  echo "==> @eat-pray-ai/${pkg_dir}@${VERSION}"
  mkdir -p "$target_dir/bin"
  cp "$src_path" "$target_dir/bin/$target_bin"
  chmod +x "$target_dir/bin/$target_bin"
  cp "$REPO_ROOT/LICENSE" "$target_dir/LICENSE"
  stamp_version "$target_dir/package.json"

  (cd "$target_dir" && npm publish "${PUBLISH_FLAGS[@]}")

  rm -rf "${target_dir:?}/bin" "${target_dir:?}/LICENSE"
  restore_version "$target_dir/package.json"
}

echo "📦 Publishing yutu npm packages v${VERSION} (tag: ${NPM_TAG})"

for pkg_dir in "${!PLATFORM_MAP[@]}"; do
  publish_platform "$pkg_dir"
done

echo "==> @eat-pray-ai/yutu@${VERSION} (root)"
root_dir="$NPM_DIR/yutu"
cp "$REPO_ROOT/LICENSE" "$root_dir/LICENSE"
cp "$REPO_ROOT/README.md" "$root_dir/README.md"
stamp_version "$root_dir/package.json"

(cd "$root_dir" && npm publish "${PUBLISH_FLAGS[@]}")

restore_version "$root_dir/package.json"
rm -f "${root_dir:?}/LICENSE" "${root_dir:?}/README.md"

echo "✅ All packages published successfully!"
