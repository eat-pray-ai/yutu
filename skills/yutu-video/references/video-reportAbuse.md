# Video ReportAbuse Command

Report abuse on a video. Use this tool when you need to report abuse on a video.

## Usage

```bash
yutu video reportAbuse [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--comments` | `-c` | Additional comments regarding the abuse report |
| `--ids` | `-i` | IDs of the videos to report abuse on |
| `--language` | `-l` | Language that the content was viewed in |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--reasonId` | `-r` | ID of the reason for reporting abuse |
| `--secondaryReasonId` | `-s` | ID of the secondary reason for reporting abuse |

## Examples

```bash
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V --secondaryReasonId V1 --language en
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId N --comments 'Spam content'
```
