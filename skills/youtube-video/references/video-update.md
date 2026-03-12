# Video Update

Update a video. Use this skill to update a video.

## Usage

```bash
yutu video update [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--categoryId` | `-g` |  | Category of the video |
| `--description` | `-d` |  | Description of the video |
| `--embeddable` | `-E` |  | Whether the video is embeddable (default true) |
| `--id` | `-i` | Yes | ID of the video to update |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--language` | `-l` |  | Language of the video |
| `--license` | `-L` |  | youtube\|creativeCommon (default "youtube") |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--playlistId` | `-y` |  | Playlist id of the video |
| `--privacy` | `-p` |  | public\|private\|unlisted |
| `--tags` | `-a` |  | Comma separated tags |
| `--thumbnail` | `-u` |  | Path to the thumbnail file |
| `--title` | `-t` |  | Title of the video |

## Examples

```bash
# Update video title
yutu video update --id dQw4w9WgXcQ --title 'New Title'
# Update video description and privacy
yutu video update --id dQw4w9WgXcQ --description 'Updated description' --privacy public
# Update video tags and category
yutu video update --id dQw4w9WgXcQ --tags 'music,pop,2024' --categoryId 10
```
