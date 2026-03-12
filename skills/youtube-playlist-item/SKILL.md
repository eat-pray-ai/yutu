---
name: youtube-playlist-item
description: "Manage YouTube playlist items. Use this skill to list items in a playlist, add new items, update items, or remove items. Useful when working with YouTube playlist item — covers listing, creating, updating, and deleting playlist item via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete items from a playlist, delete playlist item, delete my playlist item, insert a playlist item into a playlist, insert playlist item, insert my playlist item, list playlist items, list playlist item, list my playlist item, update a playlist item, update playlist item, update my playlist item"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
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
