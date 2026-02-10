---
name: yutu-playlistImage
description: Manage YouTube playlist images using the yutu CLI. Use this skill to list, insert, update, or delete playlist images.
---

# Yutu Playlist Image

## Overview

This skill allows you to manage YouTube playlist images using the `yutu` CLI tool.

## Playlist Image Operations

### List Playlist Images

Retrieve images associated with a playlist.

**Reference:** [references/playlistImage-list.md](references/playlistImage-list.md)

**Common Tasks:**

- List images: `yutu playlistImage list --parent PLAYLIST_ID`

### Insert Playlist Image

Upload a new image for a playlist.

**Reference:** [references/playlistImage-insert.md](references/playlistImage-insert.md)

**Common Tasks:**

- Insert image: `yutu playlistImage insert --playlistId PLAYLIST_ID --file image.png`

### Update Playlist Image

Update an existing playlist image.

**Reference:** [references/playlistImage-update.md](references/playlistImage-update.md)

### Delete Playlist Image

Remove a playlist image.

**Reference:** [references/playlistImage-delete.md](references/playlistImage-delete.md)

**Common Tasks:**

- Delete image: `yutu playlistImage delete --ids IMAGE_ID`

## Resources

- [references/playlistImage-list.md](references/playlistImage-list.md): List playlist images.
- [references/playlistImage-insert.md](references/playlistImage-insert.md): Insert playlist images.
- [references/playlistImage-update.md](references/playlistImage-update.md): Update playlist images.
- [references/playlistImage-delete.md](references/playlistImage-delete.md): Delete playlist images.
