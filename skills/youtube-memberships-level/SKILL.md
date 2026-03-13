---
name: youtube-memberships-level
description: "Manage YouTube memberships levels. Use this skill to list information about channel membership levels. Useful when working with YouTube memberships level — provides commands to list memberships level via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list memberships levels, list memberships level, list my memberships level"
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

# YouTube Memberships Level

Manage YouTube memberships levels. Use this skill to list information about channel membership levels.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List memberships levels | [details](references/membershipsLevel-list.md) |

## Quick Start

```bash
# Show all memberships level commands
yutu membershipsLevel --help

# List memberships level
yutu membershipsLevel list
```
