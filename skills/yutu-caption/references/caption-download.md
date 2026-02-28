# Caption Download Command

Download caption from a video.

## Usage

```bash
yutu caption download [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--file` | `-f` | Path to save the caption file |
| `--id` | `-i` | ID of the caption to download |
| `--onBehalfOf` | `-b` | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` | ID of the content owner, for YouTube content partners |
| `--tfmt` | `-t` | sbv\|srt\|vtt |
| `--tlang` | `-l` | Translate the captions into this language |

## Examples

```bash
# Download caption in srt format
yutu caption download --id CAPTION_ID --tfmt srt --file caption.srt
```
