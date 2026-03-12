# Caption List

List captions. Use this skill to list captions of a video.

## Usage

```bash
yutu caption list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` |  | IDs of the captions to list |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--onBehalfOf` | `-b` |  | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |
| `--videoId` | `-v` |  | ID of the video |

## Examples

```bash
# List captions of a video
yutu caption list --videoId dQw4w9WgXcQ
# List captions in JSON format
yutu caption list --videoId dQw4w9WgXcQ --output json
# List specific captions by IDs
yutu caption list --ids abc123,def456 --videoId dQw4w9WgXcQ
```
