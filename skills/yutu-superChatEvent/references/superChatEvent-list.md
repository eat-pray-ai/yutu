# SuperChatEvent List Command

List Super Chat events for a channel.

## Usage

```bash
yutu superChatEvent list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--hl` | `-l` | Return rendered funding amounts in specified language |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List Super Chat events
yutu superChatEvent list
```
