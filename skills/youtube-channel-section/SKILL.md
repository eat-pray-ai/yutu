---
name: youtube-channel-section
description: "Manage YouTube channel sections. Use this skill to list or delete channel sections. Useful when working with YouTube channel section — provides commands to delete and list channel section via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete channel sections, delete channel section, delete my channel section, list channel sections, list channel section, list my channel section"
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

# YouTube Channel Section

Manage YouTube channel sections. Use this skill to list or delete channel sections.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete channel sections | [details](references/channelSection-delete.md) |
| list | List channel sections | [details](references/channelSection-list.md) |

## Quick Start

```bash
# Show all channel section commands
yutu channelSection --help

# List channel section
yutu channelSection list
```
