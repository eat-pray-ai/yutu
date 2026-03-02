---
name: yutu-playlistItem
description: Manage YouTube playlist items. Use this skill when you need to list items in a playlist, add new items, update items, or remove items.
---

# Yutu PlaylistItem

## Overview

This skill allows you to manage YouTube playlist items using the `yutu` CLI tool.

## PlaylistItem Operations

### Delete items from a playlist

Delete items from a playlist. Use this tool when you need to delete items from a playlist by IDs.

```bash
# Delete a playlist item by ID
yutu playlistItem delete --ids abc123
```

**Reference:** [references/playlistItem-delete.md](references/playlistItem-delete.md)

### Insert a playlist item into a playlist

Insert a playlist item into a playlist. Use this tool when you need to insert a playlist item into a playlist.

```bash
# Add a video to a playlist
yutu playlistItem insert --kind video --playlistId PLxxx --channelId UC_x5X --kVideoId dQw4w9WgXcQ
```

**Reference:** [references/playlistItem-insert.md](references/playlistItem-insert.md)

### List playlist items

List playlist items. Use this tool when you need to list playlist items.

```bash
# List items in a playlist
yutu playlistItem list --playlistId PLxxx
```

**Reference:** [references/playlistItem-list.md](references/playlistItem-list.md)

### Update a playlist item

Update a playlist item. Use this tool when you need to update a playlist item.

```bash
# Update playlist item title
yutu playlistItem update --id abc123 --title 'Updated Title'
```

**Reference:** [references/playlistItem-update.md](references/playlistItem-update.md)

## Resources

- [references/playlistItem-delete.md](references/playlistItem-delete.md): Detailed usage of `Delete items from a playlist`
- [references/playlistItem-insert.md](references/playlistItem-insert.md): Detailed usage of `Insert a playlist item into a playlist`
- [references/playlistItem-list.md](references/playlistItem-list.md): Detailed usage of `List playlist items`
- [references/playlistItem-update.md](references/playlistItem-update.md): Detailed usage of `Update a playlist item`
