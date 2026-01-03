#!/usr/bin/env bash
# Copyright 2025 eat-pray-ai & OpenWaygate
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# Accept yutu executable path as first parameter
YUTU_PATH="${1:-}"

# If no path provided, build yutu
if [[ -z "$YUTU_PATH" ]]; then
    echo "No yutu path provided, building yutu..."

    MOD="github.com/eat-pray-ai/yutu/cmd"
    Version="${MOD}.Version=$(git describe --tags --always --dirty)"
    Commit="${MOD}.Commit=$(git rev-parse --short HEAD)"
    CommitDate="${MOD}.CommitDate=$(git log -1 --date='format:%Y-%m-%dT%H:%M:%SZ' --pretty=%cd)"
    Os="${MOD}.Os=$(go env GOOS)"
    Arch="${MOD}.Arch=$(go env GOARCH)"
    Builder="${MOD}.Builder=${GITHUB_ACTOR:-$USER}"
    ldflags="-s -X ${Version} -X ${Commit} -X ${CommitDate} -X ${Os} -X ${Arch} -X ${Builder}"

    go mod download
    go build -ldflags "${ldflags}" -o yutu .
    trap 'rm yutu' EXIT
    chmod +x yutu

    YUTU_PATH="./yutu"
else
    # Verify the provided path exists and is executable
    if [[ ! -f "$YUTU_PATH" ]]; then
        echo "Error: yutu executable not found at: $YUTU_PATH"
        exit 1
    fi
    if [[ ! -x "$YUTU_PATH" ]]; then
        echo "Error: $YUTU_PATH is not executable"
        exit 1
    fi
    echo "Using yutu at: $YUTU_PATH"
fi

# Run tests with the yutu executable
"$YUTU_PATH" help
"$YUTU_PATH" completion
"$YUTU_PATH" version

# yutu
echo "======= auth ======="
"$YUTU_PATH" auth --help

echo "======= version ======="
"$YUTU_PATH" version --help

# yutu ai
echo "======= mcp ======="
"$YUTU_PATH" mcp --help

echo "======= agent ======="
"$YUTU_PATH" agent --help

# youtube api
echo "======= activity ======="
"$YUTU_PATH" activity --help
echo "------- list -------"
"$YUTU_PATH" activity list --help

echo "======= caption ======="
"$YUTU_PATH" caption --help
echo "------- delete -------"
"$YUTU_PATH" caption delete --help
echo "------- insert -------"
"$YUTU_PATH" caption insert --help
echo "------- list -------"
"$YUTU_PATH" caption list --help
echo "------- update -------"
"$YUTU_PATH" caption update --help
echo "------- download -------"
"$YUTU_PATH" caption download --help

echo "======= channel ======="
"$YUTU_PATH" channel --help
echo "------- list -------"
"$YUTU_PATH" channel list --help
echo "------- update -------"
"$YUTU_PATH" channel update --help

echo "======= channelBanner ======="
"$YUTU_PATH" channelBanner --help
echo "------- insert -------"
"$YUTU_PATH" channelBanner insert --help

echo "======= channelSection ======="
"$YUTU_PATH" channelSection --help
echo "------- delete -------"
"$YUTU_PATH" channelSection delete --help
echo "------- insert -------"
"$YUTU_PATH" channelSection list --help

echo "======= comment ======="
"$YUTU_PATH" comment --help
echo "------- delete -------"
"$YUTU_PATH" comment delete --help
echo "------- insert -------"
"$YUTU_PATH" comment insert --help
echo "------- list -------"
"$YUTU_PATH" comment list --help
echo "------- markAsSpam -------"
"$YUTU_PATH" comment markAsSpam --help
echo "------- setModerationStatus -------"
"$YUTU_PATH" comment setModerationStatus --help
echo "------- update -------"
"$YUTU_PATH" comment update --help

echo "======= commentThread ======="
"$YUTU_PATH" commentThread --help
echo "------- list -------"
"$YUTU_PATH" commentThread list --help
echo "------- insert -------"
"$YUTU_PATH" commentThread insert --help

echo "======= i18nLanguage ======="
"$YUTU_PATH" i18nLanguage --help
echo "------- list -------"
"$YUTU_PATH" i18nLanguage list --help

echo "======= i18nRegion ======="
"$YUTU_PATH" i18nRegion --help
echo "------- list -------"
"$YUTU_PATH" i18nRegion list --help

echo "======= liveBroadcast ======="
echo "pending implementation"

echo "======= liveChatBan ======="
echo "pending implementation"

echo "======= liveChatMessage ======="
echo "pending implementation"

echo "======= liveChatModerator ======="
echo "pending implementation"

echo "======= liveStream ======="
echo "pending implementation"

echo "======= member ======="
"$YUTU_PATH" member --help
echo "------- list -------"
"$YUTU_PATH" member list --help

echo "======= membershipsLevel ======="
"$YUTU_PATH" membershipsLevel --help
echo "------- list -------"
"$YUTU_PATH" membershipsLevel list --help

echo "======= playlist ======="
"$YUTU_PATH" playlist --help
echo "------- delete -------"
"$YUTU_PATH" playlist delete --help
echo "------- insert -------"
"$YUTU_PATH" playlist insert --help
echo "------- list -------"
"$YUTU_PATH" playlist list --help
echo "------- update -------"
"$YUTU_PATH" playlist update --help

echo "======= playlistItem ======="
"$YUTU_PATH" playlistItem --help
echo "------- delete -------"
"$YUTU_PATH" playlistItem delete --help
echo "------- insert -------"
"$YUTU_PATH" playlistItem insert --help
echo "------- list -------"
"$YUTU_PATH" playlistItem list --help
echo "------- update -------"
"$YUTU_PATH" playlistItem update --help

echo "======= playlistImage ======="
"$YUTU_PATH" playlistImage --help
echo "------- delete -------"
"$YUTU_PATH" playlistImage delete --help
echo "------- insert -------"
"$YUTU_PATH" playlistImage insert --help
echo "------- list -------"
"$YUTU_PATH" playlistImage list --help
echo "------- update -------"
"$YUTU_PATH" playlistImage update --help

echo "======= search ======="
"$YUTU_PATH" search --help
echo "------- list -------"
"$YUTU_PATH" search list --help

echo "======= subscription ======="
"$YUTU_PATH" subscription --help
echo "------- delete -------"
"$YUTU_PATH" subscription delete --help
echo "------- insert -------"
"$YUTU_PATH" subscription insert --help
echo "------- list -------"
"$YUTU_PATH" subscription list --help

echo "======= superChatEvent ======="
"$YUTU_PATH" superChatEvent --help
echo "------- list -------"
"$YUTU_PATH" superChatEvent list --help

echo "======= test ======="
echo "pending implementation"

echo "======= thirdPartyLink ======="
echo "pending implementation"

echo "======= thumbnail ======="
"$YUTU_PATH" thumbnail --help
echo "------- set -------"
"$YUTU_PATH" thumbnail set --help

echo "======= video ======="
"$YUTU_PATH" video --help
echo "------- delete -------"
"$YUTU_PATH" video delete --help
echo "------- getRating -------"
"$YUTU_PATH" video getRating --help
echo "------- insert -------"
"$YUTU_PATH" video insert --help
echo "------- list -------"
"$YUTU_PATH" video list --help
echo "------- rate -------"
"$YUTU_PATH" video rate --help
echo "------- update -------"
"$YUTU_PATH" video update --help
echo "------- reportAbuse -------"
"$YUTU_PATH" video reportAbuse --help

echo "======= videoAbuseReportReason ======="
"$YUTU_PATH" videoAbuseReportReason --help
echo "------- list -------"
"$YUTU_PATH" videoAbuseReportReason list --help

echo "======= videoCategory ======="
"$YUTU_PATH" videoCategory --help
echo "------- list -------"
"$YUTU_PATH" videoCategory list --help

echo "======= watermark ======="
"$YUTU_PATH" watermark --help
echo "------- set -------"
"$YUTU_PATH" watermark set --help
echo "------- unset -------"
"$YUTU_PATH" watermark unset --help
