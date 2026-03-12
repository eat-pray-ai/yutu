---
name: youtube-channel
description: "Manage YouTube channels. Use this skill to list or update channels. Useful when working with YouTube channel — provides commands to list and update channel via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list channel information, list channel, list my channel, update channel information, update channel, update my channel"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
  required_config_paths:
    - client_secret.json
    - youtube.token.json
  env:
    - YUTU_CREDENTIAL
    - YUTU_CACHE_TOKEN
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
