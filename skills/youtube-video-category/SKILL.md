---
name: youtube-video-category
description: "Manage YouTube video categories. Use this skill to list available video categories. Useful when working with YouTube video category — provides commands to list video category via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list video categories, list video category, list my video category"
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

# YouTube Video Category

Manage YouTube video categories. Use this skill to list available video categories.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List video categories | [details](references/videoCategory-list.md) |

## Quick Start

```bash
# Show all video category commands
yutu videoCategory --help

# List video category
yutu videoCategory list
```
