#!/usr/bin/env bash

set -e

MOD="github.com/eat-pray-ai/yutu/cmd"
Version="${MOD}.Version=$(git describe --tags --always --dirty | cut -c2-)"
Commit="${MOD}.Commit=$(git rev-parse --short HEAD)"
CommitDate="${MOD}.CommitDate=$(git log -1 --date='format:%Y-%m-%d/%H.%M.%S' --pretty=%cd)"
Os="${MOD}.Os=linux"
Arch="${MOD}.Arch=$(go env GOARCH)"
ldflags="-s -X ${Version} -X ${Commit} -X ${CommitDate} -X ${Os} -X ${Arch}"

go mod download
go build -ldflags "${ldflags}" -o yutu .
trap 'rm yutu' EXIT
chmod +x yutu


./yutu help
./yutu completion

echo "======= video ======="
./yutu video --help
echo "------- delete -------"
./yutu video delete --help
echo "------- getRating -------"
./yutu video getRating --help
echo "------- insert -------"
./yutu video insert --help
echo "------- list -------"
./yutu video list --help
echo "------- rate -------"
./yutu video rate --help
echo "------- update -------"
./yutu video update --help
echo "------- reportAbuse -------"
./yutu video reportAbuse --help

echo "======= channel ======="
./yutu channel --help
echo "------- list -------"
./yutu channel list --help
echo "------- update -------"
./yutu channel update --help

echo "======= playlist ======="
./yutu playlist --help
echo "------- delete -------"
./yutu playlist delete --help
echo "------- insert -------"
./yutu playlist insert --help
echo "------- list -------"
./yutu playlist list --help
echo "------- update -------"
./yutu playlist update --help

echo "======= playlistItem ======="
./yutu playlistItem --help
echo "------- delete -------"
./yutu playlistItem delete --help
echo "------- insert -------"
./yutu playlistItem insert --help
echo "------- list -------"
./yutu playlistItem list --help
echo "------- update -------"
./yutu playlistItem update --help

echo "======= activity ======="
./yutu activity --help
echo "------- list -------"
./yutu activity list --help

echo "======= caption ======="
./yutu caption --help
echo "------- delete -------"
./yutu caption delete --help
echo "------- insert -------"
./yutu caption insert --help
echo "------- list -------"
./yutu caption list --help
echo "------- update -------"
./yutu caption update --help
echo "------- download -------"
./yutu caption download --help

echo "======= channelBanner ======="
./yutu channelBanner --help
echo "------- insert -------"
./yutu channelBanner insert --help

echo "======= channelSection ======="
./yutu channelSection --help
echo "------- delete -------"
./yutu channelSection delete --help
echo "------- insert -------"
./yutu channelSection list --help

echo "======= comment ======="
./yutu comment --help
echo "------- delete -------"
./yutu comment delete --help
echo "------- insert -------"
./yutu comment insert --help
echo "------- list -------"
./yutu comment list --help
echo "------- markAsSpam -------"
./yutu comment markAsSpam --help
echo "------- setModerationStatus -------"
./yutu comment setModerationStatus --help
echo "------- update -------"
./yutu comment update --help

echo "======= commentThread ======="
./yutu commentThread --help
echo "------- list -------"
./yutu commentThread list --help
echo "------- insert -------"
./yutu commentThread insert --help

echo "======= i18nLanguage ======="
./yutu i18nLanguage --help
echo "------- list -------"
./yutu i18nLanguage list --help

echo "======= i18nRegion ======="
./yutu i18nRegion --help
echo "------- list -------"
./yutu i18nRegion list --help

echo "======= member ======="
./yutu member --help
echo "------- list -------"
./yutu member list --help

echo "======= membershipsLevel ======="
./yutu membershipsLevel --help
echo "------- list -------"
./yutu membershipsLevel list --help

echo "======= search ======="
./yutu search --help
echo "------- list -------"
./yutu search list --help

echo "======= subscription ======="
./yutu subscription --help
echo "------- delete -------"
./yutu subscription delete --help
echo "------- insert -------"
./yutu subscription insert --help
echo "------- list -------"
./yutu subscription list --help

echo "======= superChatEvent ======="
./yutu superChatEvent --help
echo "------- list -------"
./yutu superChatEvent list --help

echo "======= thumbnail ======="
./yutu thumbnail --help
echo "------- set -------"
./yutu thumbnail set --help

echo "======= videoAbuseReport ======="
./yutu videoAbuseReportReason --help
echo "------- list -------"
./yutu videoAbuseReportReason list --help

echo "======= videoCategory ======="
./yutu videoCategory --help
echo "------- list -------"
./yutu videoCategory list --help

echo "======= watermark ======="
./yutu watermark --help
echo "------- set -------"
./yutu watermark set --help
echo "------- unset -------"
./yutu watermark unset --help

echo "======= version ======="
./yutu version --help
