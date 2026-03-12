---
name: youtube-channel-banner
description: "Manage YouTube channel banners. Use this skill to insert or upload channel banners. Useful when working with YouTube channel banner — covers listing, creating, updating, and deleting channel banner via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: insert a channel banner, insert channel banner, insert my channel banner"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Channel Banner

Manage YouTube channel banners. Use this skill to insert or upload channel banners.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| insert | Insert a channel banner | [details](references/channelBanner-insert.md) |

## Quick Start

```bash
# Show all channel banner commands
yutu channelBanner --help
```
