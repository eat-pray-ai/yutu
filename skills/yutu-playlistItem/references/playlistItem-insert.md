# Playlist Item Insert

Insert a playlist item into a playlist. Use this skill to insert a playlist item into a playlist.

## Usage

```bash
yutu playlistItem insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` | Yes | ID that YouTube uses to uniquely identify the user that added the item to the playlist |
| `--description` | `-d` |  | Description of the playlist item |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--kChannelId` | `-C` |  | ID of the channel if kind is channel |
| `--kPlaylistId` | `-Y` |  | ID of the playlist if kind is playlist |
| `--kVideoId` | `-V` |  | ID of the video if kind is video |
| `--kind` | `-k` | Yes | video\|channel\|playlist |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--playlistId` | `-y` | Yes | The id that YouTube uses to uniquely identify the playlist that the item is in |
| `--privacy` | `-p` |  | public\|private\|unlisted |
| `--title` | `-t` |  | Title of the playlist item |

## Examples

```bash
# Add a video to a playlist
yutu playlistItem insert --kind video --playlistId PLxxx --channelId UC_x5X --kVideoId dQw4w9WgXcQ
# Add a video to a playlist with privacy setting
yutu playlistItem insert --kind video --playlistId PLxxx --channelId UC_x5X --kVideoId dQw4w9 --privacy public
```
