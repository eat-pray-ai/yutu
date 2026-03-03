# PlaylistImage Delete Command

Delete playlist images.

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
# Delete a playlist image by ID
yutu playlistImage delete --ids abc123
# Delete multiple playlist images
yutu playlistImage delete --ids abc123,def456
```
