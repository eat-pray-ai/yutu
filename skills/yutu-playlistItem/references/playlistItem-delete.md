# PlaylistItem Delete Command

Delete items from a playlist by ids.

## Usage

```bash
yutu playlistItem delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the playlist items to delete |
| `--onBehalfOfContentOwner` | `-b` | |

## Examples

```bash
# Remove an item from a playlist
yutu playlistItem delete --ids ITEM_ID
```
