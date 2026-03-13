---
name: youtube-channel-banner
description: "Manage YouTube channel banners. Use this skill to insert or upload channel banners. Useful when working with YouTube channel banner — provides commands to insert channel banner via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: insert a channel banner, insert channel banner, insert my channel banner"
metadata:
  openclaw:
    requires:
      env:
        - YUTU_CREDENTIAL
        - YUTU_CACHE_TOKEN
      bins:
        - yutu
      config:
        - client_secret.json
        - youtube.token.json
    primaryEnv: YUTU_CREDENTIAL
    emoji: "\U0001F3AC\U0001F430"
    homepage: https://github.com/eat-pray-ai/yutu
    install:
      - kind: node
        package: "@eat-pray-ai/yutu"
        bins: [yutu]
---

# YouTube Channel Banner

Manage YouTube channel banners. Use this skill to insert or upload channel banners.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| insert | Insert a channel banner | [details](references/channelBanner-insert.md) |

## Quick Start

```bash
# Show all channel banner commands
yutu channelBanner --help
```
