# Playlist Update Command

Update an existing playlist.

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
yutu playlist update --id PLAYLIST_ID --title "New Title"
```
