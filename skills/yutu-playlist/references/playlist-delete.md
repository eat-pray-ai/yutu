# Playlist Delete Command

Delete playlists. Use this tool when you need to delete playlists by IDs.

## Usage

```bash
yutu playlist delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the playlists to delete |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |

## Examples

```bash
yutu playlist delete --ids PLxxxx
yutu playlist delete --ids PLxxx1,PLxxx2
```
