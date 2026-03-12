---
name: yutu-search
description: "Manage YouTube search. Use this skill to search for videos, channels, playlists, and other resources. Always use this skill when the user mentions search or wants to perform any operation on YouTube search, even if they don't explicitly ask for search management. Includes setup and installation instructions for first-time users. Triggers: search resources, list search, list my search"
---

# Yutu Search

Manage YouTube search. Use this skill to search for videos, channels, playlists, and other resources.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | Search resources | [details](references/search-list.md) |

## Quick Start

```bash
# Show all search commands
yutu search --help

# List search
yutu search list
```
