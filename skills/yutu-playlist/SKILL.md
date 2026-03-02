---
name: yutu-playlist
description: Manage YouTube playlists. Use this skill when you need to list, create, update, or delete playlists.
---

# Yutu Playlist

## Overview

This skill allows you to manage YouTube playlists using the `yutu` CLI tool.

## Playlist Operations

### Delete playlists

Delete playlists. Use this tool when you need to delete playlists by IDs.

```bash
# Delete a playlist by ID
yutu playlist delete --ids PLxxxx
```

**Reference:** [references/playlist-delete.md](references/playlist-delete.md)

### Create a new playlist

Create a new playlist. Use this tool when you need to create a new playlist.

```bash
# Create a public playlist
yutu playlist insert --title 'My Playlist' --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --privacy public
```

**Reference:** [references/playlist-insert.md](references/playlist-insert.md)

### List playlist information

List playlist information. Use this tool when you need to list playlist information.

```bash
# List my playlists
yutu playlist list --mine
```

**Reference:** [references/playlist-list.md](references/playlist-list.md)

### Update a playlist

Update a playlist. Use this tool when you need to update a playlist.

```bash
# Update playlist title
yutu playlist update --id PLxxx --title 'Updated Title'
```

**Reference:** [references/playlist-update.md](references/playlist-update.md)

## Resources

- [references/playlist-delete.md](references/playlist-delete.md): Detailed usage of `Delete playlists`
- [references/playlist-insert.md](references/playlist-insert.md): Detailed usage of `Create a new playlist`
- [references/playlist-list.md](references/playlist-list.md): Detailed usage of `List playlist information`
- [references/playlist-update.md](references/playlist-update.md): Detailed usage of `Update a playlist`
