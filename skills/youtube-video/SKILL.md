---
name: youtube-video
description: "Manage YouTube videos. Use this skill to list, upload, update, delete, get rating, or report videos. Useful when working with YouTube video — provides commands to delete, getRating, insert, list, rate, reportAbuse, and update video via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete videos, delete video, delete my video, get video ratings, getRating video, getRating my video, upload a video, insert video, insert my video, list video information, list video, list my video, rate a video, rate video, rate my video, report abuse on a video, reportAbuse video, reportAbuse my video, update a video, update video, update my video"
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

# YouTube Video

Manage YouTube videos. Use this skill to list, upload, update, delete, get rating, or report videos.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete videos | [details](references/video-delete.md) |
| getRating | Get video ratings | [details](references/video-getRating.md) |
| insert | Upload a video | [details](references/video-insert.md) |
| list | List video information | [details](references/video-list.md) |
| rate | Rate a video | [details](references/video-rate.md) |
| reportAbuse | Report abuse on a video | [details](references/video-reportAbuse.md) |
| update | Update a video | [details](references/video-update.md) |

## Quick Start

```bash
# Show all video commands
yutu video --help

# List video
yutu video list
```
