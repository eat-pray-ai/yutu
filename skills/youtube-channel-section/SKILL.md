---
name: youtube-channel-section
description: "Manage YouTube channel sections. Use this skill to list or delete channel sections. Useful when working with YouTube channel section — provides commands to delete and list channel section via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete channel sections, delete channel section, delete my channel section, list channel sections, list channel section, list my channel section"
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
