![Yutu](./asset/yutu.svg)

# yutu

[![Static Badge](https://img.shields.io/badge/gitmoji-%F0%9F%98%BF%F0%9F%90%B0%F0%9F%90%A7%E2%9D%A4%EF%B8%8F%E2%80%8D%F0%9F%A9%B9-love?style=flat-square&labelColor=%23EDD1CC&color=%23FF919F)](https://gitmoji.dev)
![Go Report Card](https://goreportcard.com/badge/github.com/eat-pray-ai/yutu?style=flat-square)
![GitHub License](https://img.shields.io/github/license/eat-pray-ai/yutu?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/eat-pray-ai/yutu/pkg/yutuber?style=flat-square)](https://pkg.go.dev/github.com/eat-pray-ai/yutu/pkg/yutuber)

![GitHub Release](https://img.shields.io/github/v/release/eat-pray-ai/yutu?sort=semver&style=flat-square&logo=go)
![GitHub Downloads](https://img.shields.io/github/downloads/eat-pray-ai/yutu/total?style=flat-square)
![GitHub Actions build Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/go-ossf-slsa3-publish.yml?style=flat-square&logo=githubactions)
![GitHub Actions CodeQL Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/codeql.yml?style=flat-square&logo=githubactions&label=CodeQL)
![GitHub Actions test Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/test.yml?style=flat-square&logo=githubactions&label=test)


yutu is a fully functional CLI for YouTube.

## Prerequisites

Before you begin, an account on [Google Cloud Platform](https://console.cloud.google.com/) is required to create a **Project** and enable these APIs for this project, in `APIs & Services -> Enable APIs and services -> + ENABLE APIS AND SERVICES`

- [YouTube Data API v3(Required)](https://console.cloud.google.com/apis/api/youtubeanalytics.googleapis.com/overview)
- [YouTube Analytics API(Optional)](https://console.cloud.google.com/apis/api/youtubeanalytics.googleapis.com/overview)
- [YouTube Reporting API(Optional)](https://console.cloud.google.com/apis/api/youtubereporting.googleapis.com/overview)

After enabling the APIs, create an `OAuth content screen` with yourself as test user, then create an `OAuth Client ID` of type `Web Application` with `http://localhost:8216` as the redirect URI.

Download this credential to your local machine with name `client_secret.json`, it should look like

```json
{
  "web": {
    "client_id": "11181119.apps.googleusercontent.com",
    "project_id": "yutu-11181119",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "XXXXXXXXXXXXXXXX",
    "redirect_uris": [
      "http://localhost:8216"
    ]
  }
}
```

To verify this credential, run the following command

```shell
❯ yutu auth --credentials client_secret.json
```

A browser window will open asking for your permission to access your YouTube account, after granting the permission, a token will be generated and saved to `youtube.token.json`.

```json
{
  "access_token": "ya29.XXXXXXXXX",
  "token_type":"Bearer",
  "refresh_token":"1//XXXXXXXXXX",
  "expiry":"2024-05-26T18:49:56.1911165+08:00"
}
```

## Installation

You can download yutu from [releases page](https://github.com/eat-pray-ai/yutu/releases/latest) directly, or use the following methods as you prefer.

### GitHub Actions

There are two actions available for yutu, one is for general purpose and the other is for uploading video to YouTube. Refer to [general](./actions/general/README.md) and [upload](./actions/upload/README.md) for more information.

### Docker

```shell
❯ docker pull ghcr.io/eat-pray-ai/yutu:latest
❯ docker run --rm ghcr.io/eat-pray-ai/yutu:latest
```

### Gopher

```shell
❯ go install https://github.com/eat-pray-ai/yutu@latest
```

### Linux

```shell
❯ curl -sSfL https://github.com/eat-pray-ai/yutu/releases/latest/download/yutu-linux-$(uname -m) -o /usr/local/bin/yutu
❯ chmod +x /usr/local/bin/yutu

```

### macOS

Homebrew is not available since this repository is not notable enough, star this repository to make it available on Homebrew.

```shell
❯ curl -sSfL https://github.com/eat-pray-ai/yutu/releases/latest/download/yutu-darwin-$(uname -m) -o /usr/local/bin/yutu
❯ chmod +x /usr/local/bin/yutu
```

### Windows

```shell
❯ winget install yutu
```

## Usage

```shell
❯ yutu help
yutu is a fully functional CLI for YouTube, which can be used to manupulate YouTube videos, playlists, channels, etc.

Usage:
  yutu [flags]
  yutu [command]

Available Commands:
  activity               List YouTube activities
  auth                   Authenticate with YouTube API
  caption                Manipulate YouTube captions
  channel                Manipulate YouTube channels
  channelBanner          Insert Youtube channelBanner
  channelSection         Manipulate channel section
  comment                Manipulate YouTube comments
  commentThread          Manipulate YouTube comment threads
  completion             Generate the autocompletion script for the specified shell
  help                   Help about any command
  i18nLanguage           List YouTube i18nLanguages
  i18nRegion             List YouTube i18nRegions
  member                 List YouTube members
  membershipsLevel       List YouTube memberships levels
  playlist               Manipulate YouTube playlists
  playlistItem           Manipulate YouTube playlist items
  search                 Search for Youtube resources
  subscription           Manipulate YouTube subscriptions
  thumbnail              Set thumbnail for a video
  version                Show the version of yutu
  video                  Manipulate YouTube videos
  videoAbuseReportReason List YouTube video abuse report reasons
  videoCategory          List YouTube video categories
  watermark              Manipulate Youtube watermarks

Flags:
  -h, --help   help for yutu

Use "yutu [command] --help" for more information about a command.
```

## Features

Please refer to [features.md](./features.md) for more information.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=eat-pray-ai/yutu&type=Date)](https://star-history.com/#eat-pray-ai/yutu&Date)
