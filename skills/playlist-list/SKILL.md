---
name: playlist-list
description: List playlists for a channel or the authenticated user.
---

# Playlist List

This skill provides instructions for listing YouTube playlists using the `yutu` CLI.

## Usage

```bash
yutu playlist list [flags]
```

## Options

- `--mine`, `-M`: Return playlists owned by the authenticated user.
- `--channelId`, `-c`: Return playlists owned by the specified channel ID.
- `--ids`, `-i`: Return playlists with the given IDs.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--hl`, `-l`: Language code for localized content.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--onBehalfOfContentOwnerChannel`, `-B`: Content owner channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet,status]`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**List my playlists:**

```bash
yutu playlist list --mine true
```

**List playlists for a channel:**

```bash
yutu playlist list --channelId CHANNEL_ID --mine false
```
