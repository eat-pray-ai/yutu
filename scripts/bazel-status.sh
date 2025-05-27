#!/usr/bin/env bash

GIT_SHORT_SHA=$(git rev-parse --short HEAD)
echo STABLE_GIT_SHORT_SHA "$GIT_SHORT_SHA"
echo STABLE_GIT_COMMIT_DATE "$(git log -1 --format=%cd --date=iso)"

GIT_TAG=$(git describe --tags --abbrev=0)
GIT_TAG_DIRTY=$(git describe --tags --dirty --abbrev=0)
# if the tag is dirty, we need to append the short SHA
if [[ "$GIT_TAG_DIRTY" != "$GIT_TAG" ]]; then
  echo STABLE_VERSION "$GIT_TAG_DIRTY-$GIT_SHORT_SHA"
else
  echo STABLE_VERSION "$GIT_TAG"
fi

echo STABLE_OS "$(uname -s | tr '[:upper:]' '[:lower:]')"
echo STABLE_ARCH "$(uname -m | sed 's/x86_64/amd64/' | tr '[:upper:]' '[:lower:]')"
