---
name: yutu-playlistItem
description: Manage YouTube playlist items using the yutu CLI. Use this skill to list items in a playlist, add new items, update items, or remove items.
---

# Yutu PlaylistItem

## Overview

This skill allows you to manage YouTube playlist items using the `yutu` CLI tool. You can add videos to playlists, list content of playlists, and remove items.

## PlaylistItem Operations

### List Playlist Items

Retrieve items within a playlist.

**Reference:** [references/playlistItem-list.md](references/playlistItem-list.md)

**Common Tasks:**

- List items: `yutu playlistItem list --playlistId PLAYLIST_ID`

### Insert/Add Item

Add a video, channel, or playlist to a playlist.

**Reference:** [references/playlistItem-insert.md](references/playlistItem-insert.md)

**Common Tasks:**

- Add video: `yutu playlistItem insert --playlistId PLAYLIST_ID --kind video --kVideoId VIDEO_ID`

### Update Item

Update metadata of a playlist item.

**Reference:** [references/playlistItem-update.md](references/playlistItem-update.md)

### Delete/Remove Item

Remove an item from a playlist.

**Reference:** [references/playlistItem-delete.md](references/playlistItem-delete.md)

**Common Tasks:**

- Remove item: `yutu playlistItem delete --ids ITEM_ID`

## Resources

- [references/playlistItem-list.md](references/playlistItem-list.md): List playlist items.
- [references/playlistItem-insert.md](references/playlistItem-insert.md): Insert playlist items.
- [references/playlistItem-update.md](references/playlistItem-update.md): Update playlist items.
- [references/playlistItem-delete.md](references/playlistItem-delete.md): Delete playlist items.
