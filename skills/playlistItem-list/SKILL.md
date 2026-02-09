---
name: playlistItem-list
description: List items in a playlist.
---

# PlaylistItem List

This skill provides instructions for listing YouTube playlist items using the `yutu` CLI.

## Usage

```bash
yutu playlistItem list [flags]
```

## Options

- `--playlistId`, `-y`: Return items within the given playlist.
- `--videoId`, `-v`: Return items associated with the given video ID.
- `--ids`, `-i`: Return items with the given IDs.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet,status]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List items in a playlist:**

```bash
yutu playlistItem list --playlistId PLAYLIST_ID
```
