# Playlist Insert Command

Create a new playlist. Use this tool when you need to create a new playlist.

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
yutu playlist insert --title 'My Playlist' --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --privacy public
yutu playlist insert --title 'Tutorial Series' --channelId UC_x5X --privacy private --description 'My tutorials'
yutu playlist insert --title 'Music' --channelId UC_x5X --privacy unlisted --tags 'music,pop'
```
