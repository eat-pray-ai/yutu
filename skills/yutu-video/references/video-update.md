# Video Update Command

Update a video. Use this tool when you need to update a video.

## Usage

```bash
yutu video update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--categoryId` | `-g` | Category of the video |
| `--description` | `-d` | Description of the video |
| `--embeddable` | `-E` | Whether the video is embeddable (default true) |
| `--id` | `-i` | ID of the video to update |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--language` | `-l` | Language of the video |
| `--license` | `-L` | youtube\|creativeCommon (default "youtube") |
| `--output` | `-o` | json\|yaml\|silent |
| `--playlistId` | `-y` | Playlist id of the video |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--tags` | `-a` | Comma separated tags |
| `--thumbnail` | `-u` | Path to the thumbnail file |
| `--title` | `-t` | Title of the video |

## Examples

```bash
yutu video update --id dQw4w9WgXcQ --title 'New Title'
yutu video update --id dQw4w9WgXcQ --description 'Updated description' --privacy public
yutu video update --id dQw4w9WgXcQ --tags 'music,pop,2024' --categoryId 10
```
