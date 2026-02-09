---
name: playlistImage-list
description: List images for a playlist.
---

# PlaylistImage List

This skill provides instructions for listing YouTube playlist images using the `yutu` CLI.

## Usage

```bash
yutu playlistImage list [flags]
```

## Options

- `--parent`, `-P`: ID of the playlist (required).
- `--maxResults`, `-n`: Maximum number of items to return.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--onBehalfOfContentOwnerChannel`, `-B`: Content owner channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,kind,snippet]`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**List images for a playlist:**

```bash
yutu playlistImage list --parent PLAYLIST_ID
```
