# PlaylistItem Update Command

Update a playlist item's info.

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
| `--onBehalfOfContentOwner` | `-b` | |
| `--output` | `-o` | json\|yaml\|silent |
| `--privacy` | `-p` | public\|private\|unlisted |
| `--title` | `-t` | Title of the playlist item |

## Examples

```bash
# Update item note (description)
yutu playlistItem update --id ITEM_ID --description "My note"
```
