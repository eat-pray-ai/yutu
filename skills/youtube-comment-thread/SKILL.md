---
name: youtube-comment-thread
description: "Manage YouTube comment threads. Use this skill to list or insert new top-level comment threads. Useful when working with YouTube comment thread — provides commands to insert and list comment thread via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: insert a new comment thread, insert comment thread, insert my comment thread, list comment threads, list comment thread, list my comment thread"
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

# YouTube Comment Thread

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
