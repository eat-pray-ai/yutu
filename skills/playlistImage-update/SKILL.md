---
name: playlistImage-update
description: Update a playlist image.
---

# PlaylistImage Update

This skill provides instructions for updating YouTube playlist images using the `yutu` CLI.

## Usage

```bash
yutu playlistImage update [flags]
```

## Options

- `--playlistId`, `-p`: ID of the playlist.
- `--type`, `-t`: Image type (e.g., 'hero').
- `--height`, `-H`: Image height.
- `--width`, `-W`: Image width.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--onBehalfOfContentOwnerChannel`, `-B`: Content owner channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update a playlist image type:**

```bash
yutu playlistImage update --playlistId PLAYLIST_ID --type hero
```
