# Activity List Command

List YouTube activities, such as likes, favorites, uploads, etc.

## Usage

```bash
yutu activity list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | ID of the channel |
| `--home` | `-H` | true\|false\|null (default true) |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` | true\|false\|null (default true) |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet,contentDetails]) |
| `--publishedAfter` | `-a` | Filter on activities published after this date |
| `--publishedBefore` | `-b` | Filter on activities published before this date |
| `--regionCode` | `-r` | Display the content as seen by viewers in this country |

## Examples

```bash
# List my activities
yutu activity list --mine

# List activities for a specific channel
yutu activity list --channelId UCxxxxxxxx
```
