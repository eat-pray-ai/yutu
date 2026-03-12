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
   - Go to `Credentials -> Create Credentials -> OAuth Client ID`, select `Web Application`
   - Add `http://localhost:8216` as an authorized redirect URI
   - Download the credential file and save it as `client_secret.json`

3. **Authenticate**:

   ```bash
   yutu auth --credential client_secret.json
   ```

   A browser window will open for you to grant YouTube access. After granting permission, a token is saved to `youtube.token.json`.

## Installation

Install `yutu` using one of these methods:


```bash
# macOS
brew install yutu

# Linux
brew install yutu

# Windows
winget install yutu

# Gopher
go install github.com/eat-pray-ai/yutu@latest
```

### Other platforms(Linux without Homebrew, etc.)

Download a prebuilt binary from the [releases page](https://github.com/eat-pray-ai/yutu/releases/latest) and place it in your PATH.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `YUTU_CREDENTIAL` | Path, base64, or JSON of OAuth client secret | `client_secret.json` |
| `YUTU_CACHE_TOKEN` | Path, base64, or JSON of cached OAuth token | `youtube.token.json` |
| `YUTU_ROOT` | Root directory for file resolution | Current working directory |
| `YUTU_LOG_LEVEL` | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR` | `INFO` |

For more details, see the [yutu README](https://github.com/eat-pray-ai/yutu#readme).
