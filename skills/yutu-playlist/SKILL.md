---
name: yutu-playlist
description: Manage YouTube playlists using the yutu CLI. Use this skill to list, create, update, or delete playlists.
---

# Yutu Playlist

## Overview

This skill allows you to manage YouTube playlists using the `yutu` CLI tool. You can list existing playlists, create new ones, update their details, and delete them.

## Playlist Operations

### List Playlists

Retrieve information about playlists.

**Reference:** [references/playlist-list.md](references/playlist-list.md)

**Common Tasks:**

- List my playlists: `yutu playlist list --mine`
- List channel playlists: `yutu playlist list --channelId CHANNEL_ID`

### Insert/Create Playlist

Create a new playlist.

**Reference:** [references/playlist-insert.md](references/playlist-insert.md)

**Common Tasks:**

- Create playlist: `yutu playlist insert --title "New Playlist" --privacy public`

### Update Playlist

Update playlist metadata.

**Reference:** [references/playlist-update.md](references/playlist-update.md)

**Common Tasks:**

- Update title: `yutu playlist update --id PLAYLIST_ID --title "New Title"`

### Delete Playlist

Remove a playlist.

**Reference:** [references/playlist-delete.md](references/playlist-delete.md)

**Common Tasks:**

- Delete playlist: `yutu playlist delete --ids PLAYLIST_ID`

## Resources

- [references/playlist-list.md](references/playlist-list.md): List playlists.
- [references/playlist-insert.md](references/playlist-insert.md): Create playlists.
- [references/playlist-update.md](references/playlist-update.md): Update playlists.
- [references/playlist-delete.md](references/playlist-delete.md): Delete playlists.
