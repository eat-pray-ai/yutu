---
name: youtube-member
description: "Manage YouTube channel members. Use this skill to list information about channel members. Useful when working with YouTube member — provides commands to list member via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list channel members, list member, list my member"
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

# YouTube Member

Manage YouTube channel members. Use this skill to list information about channel members.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List channel members | [details](references/member-list.md) |

## Quick Start

```bash
# Show all member commands
yutu member --help

# List member
yutu member list
```
