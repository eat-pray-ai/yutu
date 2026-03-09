# Video Insert

Upload a video. Use this skill to upload a video.

## Usage

```bash
yutu video insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--autoLevels` | `-A` |  | Should auto-levels be applied to the upload (default true) |
| `--categoryId` | `-g` | Yes | Category of the video |
| `--channelId` | `-c` |  | Channel id of the video |
| `--description` | `-d` |  | Description of the video |
| `--embeddable` | `-E` |  | Whether the video is embeddable (default true) |
| `--file` | `-f` | Yes | Path to the video file |
| `--forKids` | `-K` |  | Whether the video is for kids |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--language` | `-l` |  | Language of the video |
| `--license` | `-L` |  | youtube\|creativeCommon (default "youtube") |
| `--notifySubscribers` | `-N` |  | Notify the channel subscribers about the new video (default true) |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` |  | YouTube channel ID linked to the content owner |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--playlistId` | `-y` |  | Playlist id of the video |
| `--privacy` | `-p` | Yes | public\|private\|unlisted |
| `--publicStatsViewable` | `-P` |  | Whether the extended video statistics can be viewed by everyone |
| `--publishAt` | `-U` |  | Datetime when the video is scheduled to publish |
| `--stabilize` | `-S` |  | Should stabilize be applied to the upload (default true) |
| `--tags` | `-a` |  | Comma separated tags |
| `--thumbnail` | `-u` |  | Path to the thumbnail file |
| `--title` | `-t` |  | Title of the video |

## Examples

```bash
# Upload a public video
yutu video insert --file video.mp4 --title 'My Video' --categoryId 22 --privacy public
# Upload a private video with tags
yutu video insert --file video.mp4 --title 'Tutorial' --categoryId 27 --privacy private --tags 'go,tutorial'
# Upload an unlisted video with custom thumbnail
yutu video insert --file video.mp4 --title 'Music Video' --categoryId 10 --privacy unlisted --thumbnail cover.jpg
```
