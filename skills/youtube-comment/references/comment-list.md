# Comment List

List comments. Use this skill to list comments by IDs.

## Usage

```bash
yutu comment list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` |  | IDs of comments |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--maxResults` | `-n` |  | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parentId` | `-P` |  | Returns replies to the specified comment |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |
| `--textFormat` | `-t` |  | textFormatUnspecified\|html\|plainText (default "html") |

## Examples

```bash
# List replies to a comment
yutu comment list --parentId UgyXXXXXXXX --maxResults 10
# List specific comments in JSON format
yutu comment list --ids abc123,def456 --output json
# List a comment in plain text
yutu comment list --ids abc123 --textFormat plainText
```
