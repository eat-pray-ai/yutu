---
name: yutu-search
description: Manage YouTube search using the yutu CLI. Use this skill to search for videos, channels, playlists, and other resources.
---

# Yutu Search

## Overview

This skill allows you to search for YouTube resources using the `yutu` CLI tool. You can search for videos, channels, playlists, and more using various filters.

## Search Operations

### List Search Results

Search for resources matching specific criteria.

**Reference:** [references/search-list.md](references/search-list.md)

**Common Tasks:**

- Search videos: `yutu search list --q "query" --types video`
- Search channels: `yutu search list --q "query" --types channel`
- Search my content: `yutu search list --forMine`

## Resources

- [references/search-list.md](references/search-list.md): Detailed flags and usage for search operations.
