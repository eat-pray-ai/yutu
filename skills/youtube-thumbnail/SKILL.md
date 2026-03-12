---
name: youtube-thumbnail
description: "Manage YouTube video thumbnails. Use this skill to set custom thumbnails for videos. Useful when working with YouTube thumbnail — provides commands to set thumbnail via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: set a thumbnail for a video, set thumbnail, set my thumbnail"
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

# YouTube Thumbnail

Manage YouTube video thumbnails. Use this skill to set custom thumbnails for videos.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| set | Set a thumbnail for a video | [details](references/thumbnail-set.md) |

## Quick Start

```bash
# Show all thumbnail commands
yutu thumbnail --help
```
