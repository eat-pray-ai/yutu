# Thumbnail Set

Set a thumbnail for a video. Use this skill to set a thumbnail for a video.

## Usage

```bash
yutu thumbnail set [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--file` | `-f` | Yes | Path to the thumbnail file |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--videoId` | `-v` | Yes | ID of the video |

## Examples

```bash
# Set a thumbnail for a video
yutu thumbnail set --file image.jpg --videoId dQw4w9WgXcQ
# Set a thumbnail with JSON output
yutu thumbnail set --file image.png --videoId dQw4w9WgXcQ --output json
```
