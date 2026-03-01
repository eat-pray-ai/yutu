# Thumbnail Set Command

Set a thumbnail for a video. Use this tool when you need to set a thumbnail for a video.

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
yutu thumbnail set --file image.jpg --videoId dQw4w9WgXcQ
yutu thumbnail set --file image.png --videoId dQw4w9WgXcQ --output json
```
