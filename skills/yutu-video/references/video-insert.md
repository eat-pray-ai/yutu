# Video Insert Command

Upload a video to YouTube.

## Usage

```bash
yutu video insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--autoLevels` | `-A` | Should auto-levels be applied to the upload (default true) |
| `--categoryId` | `-g` | Category of the video |
| `--channelId` | `-c` | Channel id of the video |
| `--description` | `-d` | Description of the video |
| `--embeddable` | `-E` | Whether the video is embeddable (default true) |
| `--file` | `-f` | Path to the video file |
| `--forKids` | `-K` | Whether the video is for kids |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--language` | `-l` | Language of the video |
| `--license` | `-L` | youtube\|creativeCommon (default "youtube") |
| `--notifySubscribers` | `-N` | Notify the channel subscribers about the new video (default true) |
| `--onBehalfOfContentOwner` | `-b` | |
| `--onBehalfOfContentOwnerChannel` | `-B` | |
| `--output` | `-o` | json\|yaml\|silent |
| `--playlistId` | `-y` | Playlist id of the video |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--publicStatsViewable` | `-P` | Whether the extended video statistics can be viewed by everyone |
| `--publishAt` | `-U` | Datetime when the video is scheduled to publish |
| `--stabilize` | `-S` | Should stabilize be applied to the upload (default true) |
| `--tags` | `-a` | Comma separated tags |
| `--thumbnail` | `-u` | Path to the thumbnail file |
| `--title` | `-t` | Title of the video |

## Examples

```bash
# Upload a video
yutu video insert --title "My Video" --file video.mp4 --privacy public
```
