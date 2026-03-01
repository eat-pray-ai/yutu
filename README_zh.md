![Yutu](./assets/yutu.svg)

# `yutu`

[![Static Badge](https://img.shields.io/badge/gitmoji-%F0%9F%98%BF%F0%9F%90%B0%F0%9F%90%A7%E2%9D%A4%EF%B8%8F%E2%80%8D%F0%9F%A9%B9-love?style=flat-square&labelColor=%23EDD1CC&color=%23FF919F)](https://gitmoji.dev)
[![Go Report Card](https://goreportcard.com/badge/github.com/eat-pray-ai/yutu?style=flat-square)](https://goreportcard.com/report/github.com/eat-pray-ai/yutu)
[![GitHub License](https://img.shields.io/github/license/eat-pray-ai/yutu?style=flat-square)](https://github.com/eat-pray-ai/yutu?tab=Apache-2.0-1-ov-file)
[![Go Reference](https://pkg.go.dev/badge/github.com/eat-pray-ai/yutu?style=flat-square)](https://pkg.go.dev/github.com/eat-pray-ai/yutu)
[![Go Coverage](https://github.com/eat-pray-ai/yutu/wiki/coverage.svg)](https://raw.githack.com/wiki/eat-pray-ai/yutu/coverage.html)

[![GitHub Repo stars](https://img.shields.io/github/stars/eat-pray-ai/yutu?style=flat-square&logo=github)](https://github.com/eat-pray-ai/yutu/stargazers)
[![GitHub Downloads](https://img.shields.io/github/downloads/eat-pray-ai/yutu/total?style=flat-square)](https://github.com/eat-pray-ai/yutu/releases/latest)
[![GitHub Actions build Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/publish.yml?style=flat-square&logo=githubactions)](https://github.com/eat-pray-ai/yutu/actions/workflows/publish.yml)
[![GitHub Actions CodeQL Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/codeql.yml?style=flat-square&logo=githubactions&label=CodeQL)](https://github.com/eat-pray-ai/yutu/actions/workflows/codeql.yml)
[![GitHub Actions test Status](https://img.shields.io/github/actions/workflow/status/eat-pray-ai/yutu/test.yml?style=flat-square&logo=githubactions&label=test)](https://github.com/eat-pray-ai/yutu/actions/workflows/test.yml)
[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/eat-pray-ai/yutu)](https://archestra.ai/mcp-catalog/eat-pray-ai__yutu)

[![GitHub Release](https://img.shields.io/github/v/release/eat-pray-ai/yutu?sort=semver&style=flat-square&logo=go)](https://github.com/eat-pray-ai/yutu/releases/latest)
[![Homebrew Formula Version](https://img.shields.io/homebrew/v/yutu?style=flat-square&logo=homebrew)](https://formulae.brew.sh/formula/yutu)
[![WinGet Package Version](https://img.shields.io/winget/v/eat-pray-ai.yutu?style=flat-square&label=%F0%9F%93%A6%20winget
)](https://winstall.app/apps/eat-pray-ai.yutu)

[![yutu - build a fully automated YouTube Channel!](https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=988886&theme=light)](https://www.producthunt.com/posts/yutu?embed=true&utm_source=badge-featured&utm_medium=badge&utm_souce=badge-yutu)

`yutu` 是一个全功能的 MCP 服务器和 YouTube CLI 工具，用于自动化您的 YouTube 工作流程。它可以操作几乎所有的 YouTube 资源，如视频、播放列表、频道、评论、字幕等。

[![mcp demo](./assets/mcp-demo.gif)](https://asciinema.org/a/wXIHU4ciFBAKrHfaFNkMoIs12)

## 目录

- [前提条件](#前提条件)
- [安装](#安装)
  - [GitHub Actions](#github-actions)
  - [Docker](#docker)
  - [Gopher](#gopher)
  - [Linux](#linux)
  - [macOS](#macos)
  - [Windows](#windows)
  - [验证安装](#验证安装)
- [Agent](#agent)
- [MCP 服务器](#mcp-服务器)
- [使用方法](#使用方法)
- [功能特性](#功能特性)
- [贡献](#贡献)

## 前提条件

开始之前，您需要在 [Google Cloud Platform](https://console.cloud.google.com/) 上创建一个账户来新建**项目**，并为该项目启用以下 API，位置在 `APIs & Services -> Enable APIs and services -> + ENABLE APIS AND SERVICES`:

- [YouTube Data API v3（必需）](https://console.cloud.google.com/apis/api/youtubeanalytics.googleapis.com/overview)
- [YouTube Analytics API（可选）](https://console.cloud.google.com/apis/api/youtubeanalytics.googleapis.com/overview)
- [YouTube Reporting API（可选）](https://console.cloud.google.com/apis/api/youtubereporting.googleapis.com/overview)

启用 API 后，创建一个 `OAuth content screen`，将您自己设置为测试用户，然后创建一个类型为 `Web Application` 的 `OAuth Client ID`，将 `http://localhost:8216` 作为重定向 URI。

将此凭据下载到本地机器，命名为 `client_secret.json`，它应该看起来像这样：

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

要验证此凭据，请运行以下命令：

```shell
❯ yutu auth --credential client_secret.json
```

浏览器窗口将打开，要求您授权访问您的 YouTube 账户。授权后，将生成一个令牌并保存到 `youtube.token.json`。

```json
{
  "access_token": "ya29.XXXXXXXXX",
  "token_type":"Bearer",
  "refresh_token":"1//XXXXXXXXXX",
  "expiry":"2024-05-26T18:49:56.1911165+08:00"
}
```

默认情况下，`yutu` 将从当前目录读取 `client_secret.json` 和 `youtube.token.json`，`--credential/-c` 和 `--cacheToken/-t` 标志仅在 `auth` 子命令中可用。要在所有子命令中修改默认路径，请设置这些环境变量：

```shell
❯ export YUTU_CREDENTIAL=client_secret.json
❯ export YUTU_CACHE_TOKEN=youtube.token.json
# 或
❯ YUTU_CREDENTIAL=client_secret.json YUTU_CACHE_TOKEN=youtube.token.json yutu subcommand --flag value
```

### 全局环境变量

| 变量                 | 描述                                 | 默认值                  |
|--------------------|------------------------------------|----------------------|
| `YUTU_CREDENTIAL`  | OAuth 客户端密钥的路径、base64 或 JSON       | `client_secret.json` |
| `YUTU_CACHE_TOKEN` | 缓存的 OAuth 令牌的路径、base64 或 JSON      | `youtube.token.json` |
| `YUTU_ROOT`        | 文件解析的根目录                           | 当前工作目录               |
| `YUTU_LOG_LEVEL`   | 日志级别：`DEBUG`、`INFO`、`WARN`、`ERROR` | `INFO`               |

## 安装

您可以直接从[发布页面](https://github.com/eat-pray-ai/yutu/releases/latest)下载 `yutu`，或使用以下您喜欢的方法。

### GitHub Actions

yutu 有两个可用的 action，一个是通用 action，另一个专用于上传视频到 YouTube。更多信息请参考 [youtube-action](https://github.com/eat-pray-ai/youtube-action) 和 [youtube-uploader](https://github.com/eat-pray-ai/youtube-uploader)。

### Docker

```shell
❯ docker pull ghcr.io/eat-pray-ai/yutu:latest
❯ docker run --rm ghcr.io/eat-pray-ai/yutu:latest
# 确保 client_secret.json 在当前目录中
❯ docker run --rm -it -u $(id -u):$(id -g) -v $(pwd):/app -p 8216:8216 ghcr.io/eat-pray-ai/yutu:latest
```

### Gopher

```shell
❯ go install github.com/eat-pray-ai/yutu@latest
```

### Linux

```shell
❯ curl -sSfL https://raw.githubusercontent.com/eat-pray-ai/yutu/main/scripts/install.sh | bash
```

### macOS

使用 [Homebrew🍺](https://brew.sh/) 安装 `yutu`（推荐），或运行 shell 脚本。

```shell
❯ brew install yutu

# 或
❯ curl -sSfL https://raw.githubusercontent.com/eat-pray-ai/yutu/main/scripts/install.sh | bash
```

### Windows

```shell
❯ winget install yutu
```

### 验证安装

使用其关联的加密签名证明来验证 `yutu` 的完整性和来源。

```shell
# Docker
❯ gh attestation verify oci://ghcr.io/eat-pray-ai/yutu:latest --repo eat-pray-ai/yutu

# Linux 和 macOS（如果使用 shell 脚本安装）
❯ gh attestation verify $(which yutu) --repo eat-pray-ai/yutu

# Windows
❯ gh attestation verify $(where.exe yutu.exe) --repo eat-pray-ai/yutu
```

## Agent

`yutu` 提供了一种 Agent 模式来自动化 YouTube 工作流程。目前，Agent 模式仍处于实验阶段且正在积极开发中，仅支持 Google 的 Gemini 模型，需设置以下环境变量：

```shell
❯ export YUTU_ADVANCED_MODEL=google:gemini-3.1-pro-preview
❯ export YUTU_LITE_MODEL=google:gemini-3-flash-preview
❯ export YUTU_LLM_API_KEY=your_gemini_api_key
// 可选设置
❯ export GOOGLE_GEMINI_BASE_URL=https://generativelanguage.googleapis.com/
❯ export YUTU_AGENT_INSTRUCTION=Your custom instruction here
```

`YUTU_ADVANCED_MODEL` 用于编排 Agent，`YUTU_LITE_MODEL` 用于其他所有 Agent。两者均使用 `provider:modelName` 格式（目前仅支持 `google`）。如果只设置了其中一个，另一个将使用相同的值。

### Agent 环境变量

| 变量                           | 描述                                      | 必填                                               |
|------------------------------|-----------------------------------------|--------------------------------------------------|
| `YUTU_ADVANCED_MODEL`        | 编排 Agent 使用的模型（格式：`provider:modelName`） | `YUTU_ADVANCED_MODEL` 和 `YUTU_LITE_MODEL` 至少设置一个 |
| `YUTU_LITE_MODEL`            | 子 Agent 使用的模型（格式：`provider:modelName`）  | `YUTU_ADVANCED_MODEL` 和 `YUTU_LITE_MODEL` 至少设置一个 |
| `YUTU_LLM_API_KEY`           | 模型提供商的 API 密钥                           | 是                                                |
| `GOOGLE_GEMINI_BASE_URL`     | Gemini API 的 BASE_URL                   | 否                                                |
| `YUTU_AGENT_INSTRUCTION`     | 编排 Agent 的自定义指令                         | 否                                                |
| `YUTU_RETRIEVAL_INSTRUCTION` | 检索 Agent 的自定义指令                         | 否                                                |
| `YUTU_MODIFIER_INSTRUCTION`  | 修改 Agent 的自定义指令                         | 否                                                |
| `YUTU_DESTROYER_INSTRUCTION` | 删除 Agent 的自定义指令                         | 否                                                |

然后运行以下命令查看详细用法：

```
❯ yutu agent --help
❯ yutu agent --args "help"
# console mode
❯ yutu agent --args "console"
# web mode with three sub-launchers: api, a2a and webui
❯ yutu agent --args "web api a2a webui"
```

## MCP 服务器

[![在 VS Code 中安装](https://img.shields.io/badge/VS_Code-Install_Server-0098FF?style=for-the-badge&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=yutu&config=%7B%22type%22%3A%20%22stdio%22%2C%22command%22%3A%20%22yutu%22%2C%22args%22%3A%20%5B%22mcp%22%5D%2C%22env%22%3A%20%7B%22YUTU_CREDENTIAL%22%3A%20%22%2Fabsolute%2Fpath%2Fto%2Fclient_secret.json%22%2C%22YUTU_CACHE_TOKEN%22%3A%20%22%2Fabsolute%2Fpath%2Fto%2Fyoutube.token.json%22%7D%7D)
[![在 Cursor 中安装](https://cursor.com/deeplink/mcp-install-light.svg)](https://cursor.com/install-mcp?name=yutu&config=JTdCJTIyY29tbWFuZCUyMiUzQSUyMnl1dHUlMjBtY3AlMjIlMkMlMjJlbnYlMjIlM0ElN0IlMjJZVVRVX0NSRURFTlRJQUwlMjIlM0ElMjIlMkZhYnNvbHV0ZSUyRnBhdGglMkZ0byUyRmNsaWVudF9zZWNyZXQuanNvbiUyMiUyQyUyMllVVFVfQ0FDSEVfVE9LRU4lMjIlM0ElMjIlMkZhYnNvbHV0ZSUyRnBhdGglMkZ0byUyRnlvdXR1YmUudG9rZW4uanNvbiUyMiU3RCU3RA%3D%3D)

作为一个 [MCP 服务器](https://modelcontextprotocol.io/introduction)，`yutu` 可以在 MCP 客户端中使用，如 [Claude Desktop](https://modelcontextprotocol.io/quickstart/user)、[VS Code](https://code.visualstudio.com/) 或 [Cursor](https://docs.cursor.com/) 等，这允许您通过聊天的形式与 YouTube 资源进行交互。

在将 `yutu` 用作 MCP 服务器之前，请确保已安装 `yutu`（参见[安装](#安装)部分），并且您有有效的 `client_secret.json` 和 `youtube.token.json` 文件（参考[前提条件](#前提条件)部分）。

您可以通过点击上面相应的徽章将 `yutu` 添加为 VS Code 或 Cursor 中的 MCP 服务器，或手动将以下配置添加到您的 MCP 客户端。记得将 `YUTU_CREDENTIAL` 和 `YUTU_CACHE_TOKEN` 的值替换为您本地机器上的正确路径。

```json
{
  "yutu": {
    "type": "stdio",
    "command": "yutu",
    "args": [
      "mcp"
    ],
    "env": {
      "YUTU_CREDENTIAL": "/absolute/path/to/client_secret.json",
      "YUTU_CACHE_TOKEN": "/absolute/path/to/youtube.token.json"
    }
  }
}
```

## 使用方法

```shell
❯ yutu        
yutu is a fully functional MCP server and CLI for YouTube, which can manipulate almost all YouTube resources

Usage:
  yutu [flags]
  yutu [command]

Available Commands:
  activity               List YouTube activities
  auth                   Authenticate with YouTube API
  caption                Manipulate YouTube captions
  channel                Manipulate YouTube channels
  channelBanner          Insert Youtube channel banner
  channelSection         Manipulate YouTube channel sections
  comment                Manipulate YouTube comments
  commentThread          Manipulate YouTube comment threads
  completion             Generate the autocompletion script for the specified shell
  help                   Help about any command
  i18nLanguage           List YouTube i18n languages
  i18nRegion             List YouTube i18n regions
  mcp                    Start MCP server
  member                 List channel's members' info
  membershipsLevel       List memberships levels' info
  playlist               Manipulate YouTube playlists
  playlistImage          Manipulate YouTube playlist images
  playlistItem           Manipulate YouTube playlist items
  search                 Search for YouTube resources
  subscription           Manipulate YouTube subscriptions
  superChatEvent         List Super Chat events for a channel
  thumbnail              Set thumbnail for a video
  version                Show the version of yutu
  video                  Manipulate YouTube videos
  videoAbuseReportReason List YouTube video abuse report reasons
  videoCategory          List YouTube video categories
  watermark              Manipulate YouTube watermarks

Flags:
  -h, --help   help for yutu

Use "yutu [command] --help" for more information about a command.
```

## 功能特性

请参考 [FEATURES.md](docs/FEATURES.md) 获取更多信息。

## 贡献

请参考 [CONTRIBUTING.md](docs/CONTRIBUTING.md) 获取更多信息。

## Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=eat-pray-ai/yutu&type=Date)](https://star-history.com/#eat-pray-ai/yutu&Date)
