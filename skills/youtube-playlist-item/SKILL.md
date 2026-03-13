---
name: youtube-playlist-item
description: "Manage YouTube playlist items. Use this skill to list items in a playlist, add new items, update items, or remove items. Useful when working with YouTube playlist item — provides commands to delete, insert, list, and update playlist item via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete items from a playlist, delete playlist item, delete my playlist item, insert a playlist item into a playlist, insert playlist item, insert my playlist item, list playlist items, list playlist item, list my playlist item, update a playlist item, update playlist item, update my playlist item"
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

# YouTube Playlist Item

Manage YouTube playlist items. Use this skill to list items in a playlist, add new items, update items, or remove items.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete items from a playlist | [details](references/playlistItem-delete.md) |
| insert | Insert a playlist item into a playlist | [details](references/playlistItem-insert.md) |
| list | List playlist items | [details](references/playlistItem-list.md) |
| update | Update a playlist item | [details](references/playlistItem-update.md) |

## Quick Start

```bash
# Show all playlist item commands
yutu playlistItem --help

# List playlist item
yutu playlistItem list
```
