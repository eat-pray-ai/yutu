# PlaylistImage Delete Command

Delete playlist images. Use this tool when you need to delete playlist images by IDs.

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
yutu playlistImage delete --ids abc123
yutu playlistImage delete --ids abc123,def456
```
