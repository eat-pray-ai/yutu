# Video Delete

Delete videos. Use this skill to delete videos by IDs.

## Usage

```bash
yutu video delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the videos to delete |

## Examples

```bash
# Delete a video by ID
yutu video delete --ids dQw4w9WgXcQ
# Delete multiple videos
yutu video delete --ids dQw4w9WgXcQ,abc123
```
