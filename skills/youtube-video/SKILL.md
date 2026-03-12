---
name: yutu-video
description: "Manage YouTube videos. Use this skill to list, upload, update, delete, get rating, or report videos. Always use this skill when the user mentions video or wants to perform any operation on YouTube video, even if they don't explicitly ask for video management. Includes setup and installation instructions for first-time users. Triggers: delete videos, delete video, delete my video, get video ratings, getRating video, getRating my video, upload a video, insert video, insert my video, list video information, list video, list my video, rate a video, rate video, rate my video, report abuse on a video, reportAbuse video, reportAbuse my video, update a video, update video, update my video"
---

# Yutu Video

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
