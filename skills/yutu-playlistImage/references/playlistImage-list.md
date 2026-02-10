# PlaylistImage List Command

List YouTube playlist images.

## Usage

```bash
yutu playlistImage list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--onBehalfOfContentOwner` | `-b` | |
| `--onBehalfOfContentOwnerChannel` | `-B` | |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parent` | `-P` | Return PlaylistImages for this playlist id |
| `--parts` | `-p` | Comma separated parts (default [id,kind,snippet]) |

## Examples

```bash
# List images for a playlist
yutu playlistImage list --parent PLAYLIST_ID
```
