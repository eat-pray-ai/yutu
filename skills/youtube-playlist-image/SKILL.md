---
name: youtube-playlist-image
description: "Manage YouTube playlist images. Use this skill to list, insert, update, or delete playlist images. Useful when working with YouTube playlist image — provides commands to delete, insert, list, and update playlist image via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete playlist images, delete playlist image, delete my playlist image, insert a playlist image, insert playlist image, insert my playlist image, list playlist images, list playlist image, list my playlist image, update a playlist image, update playlist image, update my playlist image"
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

# YouTube Playlist Image

Manage YouTube playlist images. Use this skill to list, insert, update, or delete playlist images.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete playlist images | [details](references/playlistImage-delete.md) |
| insert | Insert a playlist image | [details](references/playlistImage-insert.md) |
| list | List playlist images | [details](references/playlistImage-list.md) |
| update | Update a playlist image | [details](references/playlistImage-update.md) |

## Quick Start

```bash
# Show all playlist image commands
yutu playlistImage --help

# List playlist image
yutu playlistImage list
```
