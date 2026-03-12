---
name: youtube-super-chat-event
description: "Manage YouTube Super Chat events. Use this skill to list Super Chat events for a channel. Useful when working with YouTube super chat event — provides commands to list super chat event via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list super chat events, list super chat event, list my super chat event"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
  required_config_paths:
    - client_secret.json
    - youtube.token.json
  env:
    - YUTU_CREDENTIAL
    - YUTU_CACHE_TOKEN
---

# YouTube Super Chat Event

Manage YouTube Super Chat events. Use this skill to list Super Chat events for a channel.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List Super Chat events | [details](references/superChatEvent-list.md) |

## Quick Start

```bash
# Show all super chat event commands
yutu superChatEvent --help

# List super chat event
yutu superChatEvent list
```
