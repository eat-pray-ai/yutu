---
name: playlistImage-delete
description: Delete playlist images.
---

# PlaylistImage Delete

This skill provides instructions for deleting YouTube playlist images using the `yutu` CLI.

## Usage

```bash
yutu playlistImage delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the playlist images to delete.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.

## Examples

**Delete a playlist image:**

```bash
yutu playlistImage delete --ids IMAGE_ID
```
