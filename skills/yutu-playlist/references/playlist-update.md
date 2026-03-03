# Playlist Update Command

Update a playlist.

## Usage

```bash
yutu playlist update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--description` | `-d` | Description of the playlist |
| `--id` | `-i` | ID of the playlist to update |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--language` | `-l` | Language of the playlist |
| `--output` | `-o` | json\|yaml\|silent |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--tags` | `-a` | Comma separated tags |
| `--title` | `-t` | Title of the playlist |

## Examples

```bash
# Update playlist title
yutu playlist update --id PLxxx --title 'Updated Title'
# Update playlist description and privacy
yutu playlist update --id PLxxx --description 'New description' --privacy public
# Update playlist tags
yutu playlist update --id PLxxx --tags 'music,pop,2024'
```
