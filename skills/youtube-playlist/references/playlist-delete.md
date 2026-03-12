# Playlist Delete

Delete playlists. Use this skill to delete playlists by IDs.

## Usage

```bash
yutu playlist delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the playlists to delete |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a playlist by ID
yutu playlist delete --ids PLxxxx
# Delete multiple playlists
yutu playlist delete --ids PLxxx1,PLxxx2
```
