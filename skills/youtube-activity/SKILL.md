---
name: youtube-activity
description: "Manage activities on YouTube. Use this skill to list channel activities. Useful when working with YouTube activity — provides commands to list activity via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list activities, list activity, list my activity"
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

# YouTube Activity

Manage activities on YouTube. Use this skill to list channel activities.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List activities | [details](references/activity-list.md) |

## Quick Start

```bash
# Show all activity commands
yutu activity --help

# List activity
yutu activity list
```
