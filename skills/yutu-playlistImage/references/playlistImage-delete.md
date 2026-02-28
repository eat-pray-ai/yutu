# PlaylistImage Delete Command

Delete YouTube playlist images by ids.

## Usage

```bash
yutu playlistImage delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the playlist images to delete |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a playlist image
yutu playlistImage delete --ids IMAGE_ID
```
