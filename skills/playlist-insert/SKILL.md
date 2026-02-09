---
name: playlist-insert
description: Create a new playlist.
---

# Playlist Insert

This skill provides instructions for creating YouTube playlists using the `yutu` CLI.

## Usage

```bash
yutu playlist insert [flags]
```

## Options

- `--title`, `-t`: Title of the playlist.
- `--description`, `-d`: Description of the playlist.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--tags`, `-a`: Comma separated tags.
- `--language`, `-l`: Language of the playlist.
- `--channelId`, `-c`: Channel ID (to create on behalf of).
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Create a public playlist:**

```bash
yutu playlist insert --title "My Playlist" --privacy public
```
