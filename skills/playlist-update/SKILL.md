---
name: playlist-update
description: Update an existing playlist.
---

# Playlist Update

This skill provides instructions for updating YouTube playlists using the `yutu` CLI.

## Usage

```bash
yutu playlist update [flags]
```

## Options

- `--id`, `-i`: ID of the playlist to update.
- `--title`, `-t`: Title of the playlist.
- `--description`, `-d`: Description of the playlist.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--tags`, `-a`: Comma separated tags.
- `--language`, `-l`: Language of the playlist.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update playlist title:**

```bash
yutu playlist update --id PLAYLIST_ID --title "New Title"
```
