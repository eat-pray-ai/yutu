# PlaylistItem Delete Command

Delete items from a playlist. Use this tool when you need to delete items from a playlist by IDs.

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
yutu playlistItem delete --ids abc123
yutu playlistItem delete --ids abc123,def456
```
