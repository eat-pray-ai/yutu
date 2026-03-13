---
name: youtube-comment
description: "Manage YouTube comments. Use this skill to list, create, update, delete, mark as spam, or set moderation status for comments. Useful when working with YouTube comment — provides commands to delete, insert, list, markAsSpam, setModerationStatus, and update comment via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete comments, delete comment, delete my comment, create a comment, insert comment, insert my comment, list comments, list comment, list my comment, mark comments as spam, markAsSpam comment, markAsSpam my comment, set comment moderation status, setModerationStatus comment, setModerationStatus my comment, update a comment on a video, update comment, update my comment"
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

# YouTube Comment

Manage YouTube comments. Use this skill to list, create, update, delete, mark as spam, or set moderation status for comments.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete comments | [details](references/comment-delete.md) |
| insert | Create a comment | [details](references/comment-insert.md) |
| list | List comments | [details](references/comment-list.md) |
| markAsSpam | Mark comments as spam | [details](references/comment-markAsSpam.md) |
| setModerationStatus | Set comment moderation status | [details](references/comment-setModerationStatus.md) |
| update | Update a comment on a video | [details](references/comment-update.md) |

## Quick Start

```bash
# Show all comment commands
yutu comment --help

# List comment
yutu comment list
```
