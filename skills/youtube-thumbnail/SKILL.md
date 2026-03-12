---
name: youtube-thumbnail
description: "Manage YouTube video thumbnails. Use this skill to set custom thumbnails for videos. Always use this skill when the user mentions thumbnail or wants to perform any operation on YouTube thumbnail, even if they don't explicitly ask for thumbnail management. Includes setup and installation instructions for first-time users. Triggers: set a thumbnail for a video, set thumbnail, set my thumbnail"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Thumbnail

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
