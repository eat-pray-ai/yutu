---
name: youtube-search
description: "Manage YouTube search. Use this skill to search for videos, channels, playlists, and other resources. Useful when working with YouTube search — provides commands to list search via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: search resources, list search, list my search"
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
