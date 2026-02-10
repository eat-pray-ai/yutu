# PlaylistItem List Command

List playlist items' info.

## Usage

```bash
yutu playlistItem list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the playlist items to list |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--onBehalfOfContentOwner` | `-b` | |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet,status]) |
| `--playlistId` | `-y` | Return the playlist items within the given playlist |
| `--videoId` | `-v` | Return the playlist items associated with the given video id |

## Examples

```bash
# List items in a playlist
yutu playlistItem list --playlistId PLAYLIST_ID

# Find playlist item for a video
yutu playlistItem list --videoId VIDEO_ID
```
