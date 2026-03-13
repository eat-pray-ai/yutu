---
name: youtube-watermark
description: "Manage YouTube watermarks. Use this skill to set or unset watermarks for channel videos. Useful when working with YouTube watermark — provides commands to set and unset watermark via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: set a watermark for channel's videos, set watermark, set my watermark, unset a watermark for channel's videos, unset watermark, unset my watermark"
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

# YouTube Watermark

Manage YouTube watermarks. Use this skill to set or unset watermarks for channel videos.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| set | Set a watermark for channel's videos | [details](references/watermark-set.md) |
| unset | Unset a watermark for channel's videos | [details](references/watermark-unset.md) |

## Quick Start

```bash
# Show all watermark commands
yutu watermark --help
```
