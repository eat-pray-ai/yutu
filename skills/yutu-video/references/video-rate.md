# Video Rate Command

Rate a video. Use this tool when you need to rate a video.

## Usage

```bash
yutu video rate [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the videos to rate |
| `--rating` | `-r` | like\|dislike\|none |

## Examples

```bash
# Like a video
yutu video rate --ids dQw4w9WgXcQ --rating like
# Dislike multiple videos
yutu video rate --ids dQw4w9WgXcQ,abc123 --rating dislike
# Remove rating from a video
yutu video rate --ids dQw4w9WgXcQ --rating none
```
