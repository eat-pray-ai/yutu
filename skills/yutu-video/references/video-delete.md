# Video Delete Command

Delete videos. Use this tool when you need to delete videos by IDs.

## Usage

```bash
yutu video delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the videos to delete |

## Examples

```bash
# Delete a video by ID
yutu video delete --ids dQw4w9WgXcQ
# Delete multiple videos
yutu video delete --ids dQw4w9WgXcQ,abc123
```
