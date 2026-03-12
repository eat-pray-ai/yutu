---
name: yutu-superChatEvent
description: "Manage YouTube Super Chat events. Use this skill to list Super Chat events for a channel. Always use this skill when the user mentions super chat event or wants to perform any operation on YouTube super chat event, even if they don't explicitly ask for super chat event management. Includes setup and installation instructions for first-time users. Triggers: list super chat events, list super chat event, list my super chat event"
---

# Yutu Super Chat Event

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
