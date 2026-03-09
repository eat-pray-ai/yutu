# Video Rate

Rate a video. Use this skill to rate a video.

## Usage

```bash
yutu video rate [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the videos to rate |
| `--rating` | `-r` | Yes | like\|dislike\|none |

## Examples

```bash
# Like a video
yutu video rate --ids dQw4w9WgXcQ --rating like
# Dislike multiple videos
yutu video rate --ids dQw4w9WgXcQ,abc123 --rating dislike
# Remove rating from a video
yutu video rate --ids dQw4w9WgXcQ --rating none
```
