# Video GetRating

Get video ratings. Use this skill to get video ratings by IDs.

## Usage

```bash
yutu video getRating [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the videos to get the rating for |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|table |

## Examples

```bash
# Get rating of a video
yutu video getRating --ids dQw4w9WgXcQ
# Get ratings of multiple videos in JSON format
yutu video getRating --ids dQw4w9WgXcQ,abc123 --output json
```
