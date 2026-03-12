# Playlist Image List

List playlist images. Use this skill to list playlist images.

## Usage

```bash
yutu playlistImage list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--maxResults` | `-n` |  | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` |  | YouTube channel ID linked to the content owner |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parent` | `-P` |  | Return PlaylistImages for this playlist id |
| `--parts` | `-p` |  | Comma separated parts (default [id,kind,snippet]) |

## Examples

```bash
# List images of a playlist
yutu playlistImage list --parent PLxxx
# List playlist images with limit in JSON format
yutu playlistImage list --parent PLxxx --maxResults 10 --output json
```
