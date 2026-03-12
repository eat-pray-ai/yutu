---
name: youtube-playlist
description: "Manage YouTube playlists. Use this skill to list, create, update, or delete playlists. Useful when working with YouTube playlist — provides commands to delete, insert, list, and update playlist via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete playlists, delete playlist, delete my playlist, create a new playlist, insert playlist, insert my playlist, list playlist information, list playlist, list my playlist, update a playlist, update playlist, update my playlist"
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

# YouTube Playlist

Manage YouTube playlists. Use this skill to list, create, update, or delete playlists.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete playlists | [details](references/playlist-delete.md) |
| insert | Create a new playlist | [details](references/playlist-insert.md) |
| list | List playlist information | [details](references/playlist-list.md) |
| update | Update a playlist | [details](references/playlist-update.md) |

## Quick Start

```bash
# Show all playlist commands
yutu playlist --help

# List playlist
yutu playlist list
```
