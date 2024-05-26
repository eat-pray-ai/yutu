![Yutu](./asset/yutu.svg)

# yutu

[![Static Badge](https://img.shields.io/badge/gitmoji-%F0%9F%98%BF%F0%9F%90%B0%F0%9F%90%A7%E2%9D%A4%EF%B8%8F%E2%80%8D%F0%9F%A9%B9-love?style=flat-square&labelColor=%23EDD1CC&color=%23FF919F)](https://gitmoji.dev)
![Go Report Card](https://goreportcard.com/badge/github.com/eat-pray-ai/yutu?style=flat-square)
![GitHub License](https://img.shields.io/github/license/eat-pray-ai/yutu?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/eat-pray-ai/yutu/pkg/yutuber?style=flat-square)](https://pkg.go.dev/github.com/eat-pray-ai/yutu/pkg/yutuber)

![GitHub Release](https://img.shields.io/github/v/release/eat-pray-ai/yutu?sort=semver&style=flat-square&logo=go)
![GitHub Downloads](https://img.shields.io/github/downloads/eat-pray-ai/yutu/total?style=flat-square)
![GitHub Actions build Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/go-ossf-slsa3-publish.yml?style=flat-square&logo=githubactions)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/codeql.yml?style=flat-square&logo=githubactions&label=CodeQL)


yutu is a fully functional CLI for YouTube.

## Usage

```shell
❯ yutu help
yutu is a fully functional CLI for YouTube, which can be used to manupulate YouTube videos, playlists, channels, etc.

Usage:
  yutu [flags]
  yutu [command]

Available Commands:
  activity               manipulate YouTube activities
  channel                manipulate YouTube channels
  channelBanner          manipulate Youtube channelBanner
  completion             Generate the autocompletion script for the specified shell
  help                   Help about any command
  i18nLanguage           manipulate YouTube i18nLanguages
  i18nRegion             manipulate YouTube i18nRegions
  member                 manipulate YouTube members
  membershipsLevel       manipulate YouTube memberships levels
  playlist               manipulate YouTube playlists
  playlistItem           manipulate YouTube playlist items
  search                 Search for youtube resources
  subscription           manipulate YouTube subscriptions
  version                Show the version of yutu
  video                  manipulate YouTube videos
  videoAbuseReportReason manipulate YouTube video abuse report reasons
  videoCategory          manipulate YouTube video categories

Flags:
  -h, --help   help for yutu

Use "yutu [command] --help" for more information about a command.
```

## Features

Here are the features that are currently supported by yutu, and the ones that are planned to be supported in the future. The quota costs for each feature is also mentioned since there is a quota limits of 10,000 units/day.

- videos
  - [x] list, 1
  - [x] insert, 1600
  - [x] update, 50
  - [x] rate, 50
  - [x] getRating, 1
  - [ ] reportAbuse, 50
  - [ ] delete, 50
- channels
  - [x] list, 1
  - [x] update, 50
- playlists
  - [x] list, 1
  - [x] insert, 50
  - [x] update, 50
  - [x] delete, 50
- playlistItems
  - [x] list, 1
  - [x] insert, 50
  - [x] update, 50
  - [x] delete, 50
- activities
  - [x] list, 1
- captions
  - [ ] list, 50
  - [ ] insert, 400
  - [ ] update, 450
  - [ ] delete, 50
- channelBanners
  - [x] insert, 50
- channelSections
  - [ ] list, 1
  - [ ] insert, 50
  - [ ] update, 50
  - [ ] delete, 50
- comments
  - [ ] list, 1
  - [ ] insert, 50
  - [ ] update, 50
  - [ ] setModerationStatus, 50
  - [ ] delete, 50
- commentThreads
  - [ ] list, 1
  - [ ] insert, 50
  - [ ] update, 50
- <s>guideCategories</s>
  - [x] <s>list, 1 deprecated API</s>
- i18nLanguages
  - [x] list, 1
- i18nRegions
  - [x] list, 1
- members
  - [x] list, 1 [🚫issue #3](https://github.com/eat-pray-ai/yutu/issues/3)
- membershipsLevels
  - [x] list, 1 [🚫issue #3](https://github.com/eat-pray-ai/yutu/issues/3)
- search
  - [x] list, 100
- subscriptions
  - [x] list, 1
  - [x] insert, 50
  - [x] delete, 50
- thumbnails
  - [x] set, 50
- videoAbuseReportReasons
  - [x] list, 1
- videoCategories
  - [x] list, 1
- watermarks
  - [ ] set, 50
  - [ ] unset, 50

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=eat-pray-ai/yutu&type=Date)](https://star-history.com/#eat-pray-ai/yutu&Date)
