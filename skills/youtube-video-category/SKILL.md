---
name: youtube-video-category
description: "Manage YouTube video categories. Use this skill to list available video categories. Useful when working with YouTube video category — covers listing, creating, updating, and deleting video category via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list video categories, list video category, list my video category"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Video Category

Manage YouTube video categories. Use this skill to list available video categories.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List video categories | [details](references/videoCategory-list.md) |

## Quick Start

```bash
# Show all video category commands
yutu videoCategory --help

# List video category
yutu videoCategory list
```
