# Watermark Set Command

Set watermark for channel's video by channel id.

## Usage

```bash
yutu watermark set [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | ID of the channel |
| `--durationMs` | `-d` | Duration in milliseconds for which the watermark should be displayed |
| `--file` | `-f` | Path to the watermark file |
| `--inVideoPosition` | `-p` | topLeft\|topRight\|bottomLeft\|bottomRight |
| `--offsetMs` | `-m` | Defines the time at which the watermark will appear |
| `--offsetType` | `-t` | offsetFromStart\|offsetFromEnd |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Set watermark for my channel
yutu watermark set --file watermark.png --inVideoPosition bottomRight
```
