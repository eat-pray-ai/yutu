#!/usr/bin/env bash
set -euo pipefail

GIT_SHORT_SHA=$(git rev-parse --short HEAD)
echo STABLE_GIT_SHORT_SHA "$GIT_SHORT_SHA"
echo STABLE_GIT_COMMIT_DATE "$(git log -1 --date='format:%Y-%m-%dT%H:%M:%SZ' --pretty=%cd)"

GIT_TAG=$(git describe --tags --abbrev=0)
GIT_TAG_DIRTY=$(git describe --tags --always --dirty)
# if the tag is dirty, we need to append the short SHA
if [[ "$GIT_TAG_DIRTY" != "$GIT_TAG" ]]; then
  echo STABLE_VERSION "$GIT_TAG_DIRTY"
else
  echo STABLE_VERSION "$GIT_TAG"
fi

echo STABLE_OS "$(uname -s | tr '[:upper:]' '[:lower:]')"
echo STABLE_ARCH "$(uname -m | sed 's/x86_64/amd64/' | tr '[:upper:]' '[:lower:]')"
