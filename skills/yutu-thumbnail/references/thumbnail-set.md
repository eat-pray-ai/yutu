# Thumbnail Set Command

Set thumbnail for a video.

## Usage

```bash
yutu thumbnail set [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--file` | `-f` | Path to the thumbnail file |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--videoId` | `-v` | ID of the video |

## Examples

```bash
# Set video thumbnail
yutu thumbnail set --videoId VIDEO_ID --file thumbnail.jpg
```
