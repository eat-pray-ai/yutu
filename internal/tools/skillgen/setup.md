<!-- Keep in sync with README.md Prerequisites and Installation sections. -->

# Yutu Setup Guide

Before using any yutu skill, you need to install `yutu` and configure YouTube API credentials.

## Prerequisites

An account on [Google Cloud Platform](https://console.cloud.google.com/) is required. Set up the following:

1. **Create a GCP Project** and enable these APIs under `APIs & Services -> Enable APIs and services`:
   - [YouTube Data API v3](https://console.cloud.google.com/apis/api/youtube.googleapis.com/overview) (Required)
   - [YouTube Analytics API](https://console.cloud.google.com/apis/api/youtubeanalytics.googleapis.com/overview) (Optional)
   - [YouTube Reporting API](https://console.cloud.google.com/apis/api/youtubereporting.googleapis.com/overview) (Optional)

2. **Create OAuth credentials**:
   - Go to `APIs & Services -> OAuth consent screen`, create a consent screen with yourself as a test user
   - Go to `Credentials -> Create Credentials -> OAuth Client ID`, select `Desktop app`
   - Download the credential file and save it as `client_secret.json`

3. **Authenticate**:

   ```bash
   yutu auth --credential client_secret.json
   ```

   A browser window will open for you to grant YouTube access. After granting permission, a token is saved to `youtube.token.json`.

## Installation

Install `yutu` using one of these methods:

```bash
# Node.js (all platforms)
npm i -g @eat-pray-ai/yutu

# macOS
brew install yutu

# Linux / macOS (shell script)
curl -sSfL https://raw.githubusercontent.com/eat-pray-ai/yutu/main/scripts/install.sh | bash

# Windows
winget install yutu
```

### Other platforms

Download a prebuilt binary from the [releases page](https://github.com/eat-pray-ai/yutu/releases/latest) and place it in your PATH.

## MCP Server

Add `yutu` as an MCP server after installation and authentication:

```bash
# Claude Code
claude mcp add -e YUTU_CREDENTIAL=/absolute/path/to/client_secret.json \
  -e YUTU_CACHE_TOKEN=/absolute/path/to/youtube.token.json \
  yutu -- yutu mcp

# Codex
codex mcp add --env YUTU_CREDENTIAL=/absolute/path/to/client_secret.json \
  --env YUTU_CACHE_TOKEN=/absolute/path/to/youtube.token.json \
  yutu -- yutu mcp
```

For VS Code, Cursor, or other tools, add to your MCP settings:

```json
{
  "yutu": {
    "type": "stdio",
    "command": "yutu",
    "args": ["mcp"],
    "env": {
      "YUTU_CREDENTIAL": "/absolute/path/to/client_secret.json",
      "YUTU_CACHE_TOKEN": "/absolute/path/to/youtube.token.json"
    }
  }
}
```

## Skills

```bash
npx skills add https://github.com/eat-pray-ai/yutu/tree/main/skills/youtube
```

## Environment Variables

| Variable | Description                                  | Default |
|----------|----------------------------------------------|---------|
| `YUTU_CREDENTIAL` | Path, Base64, or JSON of OAuth client secret | `client_secret.json` |
| `YUTU_CACHE_TOKEN` | Path, Base64, or JSON of cached OAuth token  | `youtube.token.json` |
| `YUTU_ROOT` | Root directory for file resolution           | Current working directory |
| `YUTU_LOG_LEVEL` | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR`  | `INFO` |

For more details, see the [README](https://github.com/eat-pray-ai/yutu#readme).
