---
name: youtube-thumbnail
description: "Manage YouTube video thumbnails. Use this skill to set custom thumbnails for videos. Useful when working with YouTube thumbnail — provides commands to set thumbnail via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: set a thumbnail for a video, set thumbnail, set my thumbnail"
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
