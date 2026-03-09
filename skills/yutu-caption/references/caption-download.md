# Caption Download

Download a caption. Use this skill to download a caption from a video.

## Usage

```bash
yutu caption download [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--file` | `-f` | Yes | Path to save the caption file |
| `--id` | `-i` | Yes | ID of the caption to download |
| `--onBehalfOf` | `-b` |  | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` |  | ID of the content owner, for YouTube content partners |
| `--tfmt` | `-t` |  | sbv\|srt\|vtt |
| `--tlang` | `-l` |  | Translate the captions into this language |

## Examples

```bash
# Download a caption as SRT
yutu caption download --id abc123 --file subtitle.srt
# Download a caption as VTT
yutu caption download --id abc123 --file subtitle.vtt --tfmt vtt
# Download a caption translated to French
yutu caption download --id abc123 --file subtitle.srt --tlang fr
```
