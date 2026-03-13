---
name: youtube-channel
description: "Manage YouTube channels. Use this skill to list or update channels. Useful when working with YouTube channel — provides commands to list and update channel via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list channel information, list channel, list my channel, update channel information, update channel, update my channel"
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

# YouTube Channel

Manage YouTube channels. Use this skill to list or update channels.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List channel information | [details](references/channel-list.md) |
| update | Update channel information | [details](references/channel-update.md) |

## Quick Start

```bash
# Show all channel commands
yutu channel --help

# List channel
yutu channel list
```
