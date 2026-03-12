---
name: youtube-watermark
description: "Manage YouTube watermarks. Use this skill to set or unset watermarks for channel videos. Always use this skill when the user mentions watermark or wants to perform any operation on YouTube watermark, even if they don't explicitly ask for watermark management. Includes setup and installation instructions for first-time users. Triggers: set a watermark for channel's videos, set watermark, set my watermark, unset a watermark for channel's videos, unset watermark, unset my watermark"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Watermark

Manage YouTube watermarks. Use this skill to set or unset watermarks for channel videos.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| set | Set a watermark for channel's videos | [details](references/watermark-set.md) |
| unset | Unset a watermark for channel's videos | [details](references/watermark-unset.md) |

## Quick Start

```bash
# Show all watermark commands
yutu watermark --help
```
