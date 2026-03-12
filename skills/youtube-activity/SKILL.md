---
name: youtube-activity
description: "Manage activities on YouTube. Use this skill to list channel activities. Always use this skill when the user mentions activity or wants to perform any operation on YouTube activity, even if they don't explicitly ask for activity management. Includes setup and installation instructions for first-time users. Triggers: list activities, list activity, list my activity"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Activity

Manage activities on YouTube. Use this skill to list channel activities.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List activities | [details](references/activity-list.md) |

## Quick Start

```bash
# Show all activity commands
yutu activity --help

# List activity
yutu activity list
```
