---
name: playlistImage-insert
description: Insert (upload) a playlist image.
---

# PlaylistImage Insert

This skill provides instructions for inserting YouTube playlist images using the `yutu` CLI.

## Usage

```bash
yutu playlistImage insert [flags]
```

## Options

- `--playlistId`, `-p`: ID of the playlist.
- `--file`, `-f`: Path to the image file.
- `--type`, `-t`: Image type (e.g., 'hero').
- `--height`, `-H`: Image height.
- `--width`, `-W`: Image width.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--onBehalfOfContentOwnerChannel`, `-B`: Content owner channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Upload a playlist hero image:**

```bash
yutu playlistImage insert --playlistId PLAYLIST_ID --file image.jpg --type hero
```
