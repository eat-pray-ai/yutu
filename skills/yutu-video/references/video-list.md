# Video List Command

List video's info.

## Usage

```bash
yutu video list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--chart` | `-c` | chartUnspecified\|mostPopular |
| `--hl` | `-l` | Specifies the localization language |
| `--ids` | `-i` | Return videos with the given ids |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--locale` | `-L` | |
| `--maxHeight` | `-H` | |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--maxWidth` | `-W` | |
| `--myRating` | `-R` | Return videos liked/disliked by the authenticated user |
| `--onBehalfOfContentOwner` | `-b` | |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet,status]) |
| `--regionCode` | `-r` | Specific to the specified region |
| `--videoCategoryId` | `-g` | Category of the video |

## Examples

```bash
# List my videos
yutu video list --myRating like

# List most popular videos
yutu video list --chart mostPopular
```
