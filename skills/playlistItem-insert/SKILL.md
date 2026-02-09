---
name: playlistItem-insert
description: Insert a new item (video, channel, playlist) into a playlist.
---

# PlaylistItem Insert

This skill provides instructions for inserting YouTube playlist items using the `yutu` CLI.

## Usage

```bash
yutu playlistItem insert [flags]
```

## Options

- `--playlistId`, `-y`: ID of the playlist to add item to.
- `--kind`, `-k`: `video`, `channel`, `playlist` (resource being added).
- `--kVideoId`, `-V`: ID of the video (if kind is video).
- `--kChannelId`, `-C`: ID of the channel (if kind is channel).
- `--kPlaylistId`, `-Y`: ID of the playlist (if kind is playlist).
- `--title`, `-t`: Title of the playlist item.
- `--description`, `-d`: Description of the playlist item.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--channelId`, `-c`: User ID adding the item.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Add a video to a playlist:**

```bash
yutu playlistItem insert --playlistId PLAYLIST_ID --kind video --kVideoId VIDEO_ID
```
