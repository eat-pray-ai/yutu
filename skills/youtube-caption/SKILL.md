---
name: youtube-caption
description: "Manage YouTube video captions. Use this skill to list, insert, update, download, or delete video captions. Useful when working with YouTube caption — provides commands to delete, download, insert, list, and update caption via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete captions, delete caption, delete my caption, download a caption, download caption, download my caption, insert a caption, insert caption, insert my caption, list captions, list caption, list my caption, update a video caption, update caption, update my caption"
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

# YouTube Caption

Manage YouTube video captions. Use this skill to list, insert, update, download, or delete video captions.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete captions | [details](references/caption-delete.md) |
| download | Download a caption | [details](references/caption-download.md) |
| insert | Insert a caption | [details](references/caption-insert.md) |
| list | List captions | [details](references/caption-list.md) |
| update | Update a video caption | [details](references/caption-update.md) |

## Quick Start

```bash
# Show all caption commands
yutu caption --help

# List caption
yutu caption list
```
