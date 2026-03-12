# Playlist Insert

Create a new playlist. Use this skill to create a new playlist.

## Usage

```bash
yutu playlist insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` | Yes | Channel id of the playlist |
| `--description` | `-d` |  | Description of the playlist |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--language` | `-l` |  | Language of the playlist |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--privacy` | `-p` | Yes | public\|private\|unlisted |
| `--tags` | `-a` |  | Comma separated tags |
| `--title` | `-t` | Yes | Title of the playlist |

## Examples

```bash
# Create a public playlist
yutu playlist insert --title 'My Playlist' --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --privacy public
# Create a private playlist with description
yutu playlist insert --title 'Tutorial Series' --channelId UC_x5X --privacy private --description 'My tutorials'
# Create an unlisted playlist with tags
yutu playlist insert --title 'Music' --channelId UC_x5X --privacy unlisted --tags 'music,pop'
```
