# Playlist Item Delete

Delete items from a playlist. Use this skill to delete items from a playlist by IDs.

## Usage

```bash
yutu playlistItem delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the playlist items to delete |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a playlist item by ID
yutu playlistItem delete --ids abc123
# Delete multiple playlist items
yutu playlistItem delete --ids abc123,def456
```
