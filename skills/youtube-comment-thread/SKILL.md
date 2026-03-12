---
name: youtube-comment-thread
description: "Manage YouTube comment threads. Use this skill to list or insert new top-level comment threads. Always use this skill when the user mentions comment thread or wants to perform any operation on YouTube comment thread, even if they don't explicitly ask for comment thread management. Includes setup and installation instructions for first-time users. Triggers: insert a new comment thread, insert comment thread, insert my comment thread, list comment threads, list comment thread, list my comment thread"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Comment Thread

Manage YouTube comment threads. Use this skill to list or insert new top-level comment threads.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| insert | Insert a new comment thread | [details](references/commentThread-insert.md) |
| list | List comment threads | [details](references/commentThread-list.md) |

## Quick Start

```bash
# Show all comment thread commands
yutu commentThread --help

# List comment thread
yutu commentThread list
```
