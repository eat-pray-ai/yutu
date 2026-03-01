# PlaylistItem Update Command

Update a playlist item. Use this tool when you need to update a playlist item.

## Usage

```bash
yutu playlistItem update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--description` | `-d` | Description of the playlist item |
| `--id` | `-i` | ID of the playlist item to update |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--output` | `-o` | json\|yaml\|silent |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--title` | `-t` | Title of the playlist item |

## Examples

```bash
yutu playlistItem update --id abc123 --title 'Updated Title'
yutu playlistItem update --id abc123 --description 'New description' --privacy public
yutu playlistItem update --id abc123 --privacy private --output json
```
