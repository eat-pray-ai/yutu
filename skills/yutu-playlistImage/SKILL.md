---
name: yutu-playlistImage
description: Manage YouTube playlist images. Use this skill when you need to list, insert, update, or delete playlist images.
---

# Yutu PlaylistImage

## Overview

This skill allows you to manage YouTube playlist images using the `yutu` CLI tool.

## PlaylistImage Operations

### Delete playlist images

Delete playlist images. Use this tool when you need to delete playlist images by IDs.

```bash
# Delete a playlist image by ID
yutu playlistImage delete --ids abc123
```

**Reference:** [references/playlistImage-delete.md](references/playlistImage-delete.md)

### Insert a playlist image

Insert a playlist image. Use this tool when you need to insert a YouTube playlist image for a given playlist ID.

```bash
# Insert a playlist cover image
yutu playlistImage insert --file cover.jpg --playlistId PLxxx
```

**Reference:** [references/playlistImage-insert.md](references/playlistImage-insert.md)

### List playlist images

List playlist images. Use this tool when you need to list playlist images.

```bash
# List images of a playlist
yutu playlistImage list --parent PLxxx
```

**Reference:** [references/playlistImage-list.md](references/playlistImage-list.md)

### Update a playlist image

Update a playlist image. Use this tool when you need to update a playlist image.

```bash
# Update a playlist image
yutu playlistImage update --playlistId PLxxx --type hero --width 2048 --height 1152
```

**Reference:** [references/playlistImage-update.md](references/playlistImage-update.md)

## Resources

- [references/playlistImage-delete.md](references/playlistImage-delete.md): Detailed usage of `Delete playlist images`
- [references/playlistImage-insert.md](references/playlistImage-insert.md): Detailed usage of `Insert a playlist image`
- [references/playlistImage-list.md](references/playlistImage-list.md): Detailed usage of `List playlist images`
- [references/playlistImage-update.md](references/playlistImage-update.md): Detailed usage of `Update a playlist image`
