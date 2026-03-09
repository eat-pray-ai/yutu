# Activity List

List activities. Use this skill to list activities such as uploads, likes, and favorites.

## Usage

```bash
yutu activity list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` |  | ID of the channel |
| `--home` | `-H` |  | true\|false\|null (default true) |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--maxResults` | `-n` |  | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` |  | true\|false\|null (default true) |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet,contentDetails]) |
| `--publishedAfter` | `-a` |  | Filter on activities published after this date |
| `--publishedBefore` | `-b` |  | Filter on activities published before this date |
| `--regionCode` | `-r` |  | Display the content as seen by viewers in this country |

## Examples

```bash
# List my activities
yutu activity list --mine
# List activities by channel ID with limit
yutu activity list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10
# List activities after a date in JSON format
yutu activity list --publishedAfter 2024-01-01T00:00:00Z --output json
```
