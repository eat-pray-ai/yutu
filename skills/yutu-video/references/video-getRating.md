# Video GetRating Command

Get video ratings. Use this tool when you need to get video ratings by IDs.

## Usage

```bash
yutu video getRating [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the videos to get the rating for |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--output` | `-o` | json\|yaml\|table |

## Examples

```bash
yutu video getRating --ids dQw4w9WgXcQ
yutu video getRating --ids dQw4w9WgXcQ,abc123 --output json
```
