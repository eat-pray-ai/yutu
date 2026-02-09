---
name: playlistItem-delete
description: Delete items from a playlist.
---

# PlaylistItem Delete

This skill provides instructions for deleting YouTube playlist items using the `yutu` CLI.

## Usage

```bash
yutu playlistItem delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the playlist items to delete.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.

## Examples

**Delete a playlist item:**

```bash
yutu playlistItem delete --ids ITEM_ID
```
