# Video ReportAbuse

Report abuse on a video. Use this skill to report abuse on a video.

## Usage

```bash
yutu video reportAbuse [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--comments` | `-c` |  | Additional comments regarding the abuse report |
| `--ids` | `-i` | Yes | IDs of the videos to report abuse on |
| `--language` | `-l` |  | Language that the content was viewed in |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--reasonId` | `-r` | Yes | ID of the reason for reporting abuse |
| `--secondaryReasonId` | `-s` |  | ID of the secondary reason for reporting abuse |

## Examples

```bash
# Report abuse on a video
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V
# Report abuse with secondary reason and language
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V --secondaryReasonId V1 --language en
# Report abuse with comments
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId N --comments 'Spam content'
```
