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
# Update playlist item title
yutu playlistItem update --id abc123 --title 'Updated Title'
# Update playlist item description and privacy
yutu playlistItem update --id abc123 --description 'New description' --privacy public
# Update playlist item privacy with JSON output
yutu playlistItem update --id abc123 --privacy private --output json
```
