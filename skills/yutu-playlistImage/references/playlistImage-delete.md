# Playlist Image Delete

Delete playlist images. Use this skill to delete playlist images by IDs.

## Usage

```bash
yutu playlistImage delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the playlist images to delete |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a playlist image by ID
yutu playlistImage delete --ids abc123
# Delete multiple playlist images
yutu playlistImage delete --ids abc123,def456
```
