# Playlist List Command

List playlist's info.

## Usage

```bash
yutu playlist list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | Return the playlists owned by the specified channel id |
| `--hl` | `-l` | Return content in specified language |
| `--ids` | `-i` | Return the playlists with the given Ids |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` | Return the playlists owned by the authenticated user (default true) |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` | YouTube channel ID linked to the content owner |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet,status]) |

## Examples

```bash
# List my playlists
yutu playlist list --mine

# List channel's playlists
yutu playlist list --channelId CHANNEL_ID
```
