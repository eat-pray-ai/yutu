---
name: youtube-comment
description: "Manage YouTube comments. Use this skill to list, create, update, delete, mark as spam, or set moderation status for comments. Always use this skill when the user mentions comment or wants to perform any operation on YouTube comment, even if they don't explicitly ask for comment management. Includes setup and installation instructions for first-time users. Triggers: delete comments, delete comment, delete my comment, create a comment, insert comment, insert my comment, list comments, list comment, list my comment, mark comments as spam, markAsSpam comment, markAsSpam my comment, set comment moderation status, setModerationStatus comment, setModerationStatus my comment, update a comment on a video, update comment, update my comment"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Comment

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
