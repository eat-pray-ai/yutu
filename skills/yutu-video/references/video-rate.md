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
yutu video rate --ids dQw4w9WgXcQ --rating like
yutu video rate --ids dQw4w9WgXcQ,abc123 --rating dislike
yutu video rate --ids dQw4w9WgXcQ --rating none
```
