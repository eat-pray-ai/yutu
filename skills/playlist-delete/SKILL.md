---
name: playlist-delete
description: Delete a playlist by ID.
---

# Playlist Delete

This skill provides instructions for deleting YouTube playlists using the `yutu` CLI.

## Usage

```bash
yutu playlist delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the playlists to delete (comma-separated).
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.

## Examples

**Delete a playlist:**

```bash
yutu playlist delete --ids PLAYLIST_ID
```
