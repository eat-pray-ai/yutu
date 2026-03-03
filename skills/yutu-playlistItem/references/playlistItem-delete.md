# PlaylistItem Delete Command

Delete items from a playlist.

## Usage

```bash
yutu playlistItem delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the playlist items to delete |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a playlist item by ID
yutu playlistItem delete --ids abc123
# Delete multiple playlist items
yutu playlistItem delete --ids abc123,def456
```
