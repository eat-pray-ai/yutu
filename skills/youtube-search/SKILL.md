---
name: youtube-search
description: "Manage YouTube search. Use this skill to search for videos, channels, playlists, and other resources. Useful when working with YouTube search — provides commands to list search via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: search resources, list search, list my search"
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

# YouTube Search

Manage YouTube search. Use this skill to search for videos, channels, playlists, and other resources.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | Search resources | [details](references/search-list.md) |

## Quick Start

```bash
# Show all search commands
yutu search --help

# List search
yutu search list
```
