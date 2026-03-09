# Playlist Item List

List playlist items. Use this skill to list playlist items.

## Usage

```bash
yutu playlistItem list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` |  | IDs of the playlist items to list |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--maxResults` | `-n` |  | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet,status]) |
| `--playlistId` | `-y` |  | Return the playlist items within the given playlist |
| `--videoId` | `-v` |  | Return the playlist items associated with the given video id |

## Examples

```bash
# List items in a playlist
yutu playlistItem list --playlistId PLxxx
# List playlist items with limit in JSON format
yutu playlistItem list --playlistId PLxxx --maxResults 20 --output json
# List specific playlist items by IDs
yutu playlistItem list --ids abc123,def456
```
