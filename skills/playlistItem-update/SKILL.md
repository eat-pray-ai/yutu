---
name: playlistItem-update
description: Update a playlist item (e.g., position, note).
---

# PlaylistItem Update

This skill provides instructions for updating YouTube playlist items using the `yutu` CLI.

## Usage

```bash
yutu playlistItem update [flags]
```

## Options

- `--id`, `-i`: ID of the playlist item to update.
- `--title`, `-t`: Title of the playlist item.
- `--description`, `-d`: Description of the playlist item.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update a playlist item description:**

```bash
yutu playlistItem update --id ITEM_ID --description "New Description"
```
