# Playlist Insert Command

Create a new playlist.

## Usage

```bash
yutu playlist insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | Channel id of the playlist |
| `--description` | `-d` | Description of the playlist |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--language` | `-l` | Language of the playlist |
| `--output` | `-o` | json\|yaml\|silent |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--tags` | `-a` | Comma separated tags |
| `--title` | `-t` | Title of the playlist |

## Examples

```bash
# Create a public playlist
yutu playlist insert --title "My Playlist" --privacy public
```
