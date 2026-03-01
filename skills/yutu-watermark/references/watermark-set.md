# Watermark Set Command

Set a watermark for channel's videos. Use this tool when you need to set a watermark for channel's videos.

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
yutu watermark set --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file watermark.png
yutu watermark set --channelId UC_x5X --file logo.png --inVideoPosition bottomRight --offsetType offsetFromEnd --offsetMs 1000
yutu watermark set --channelId UC_x5X --file logo.png --durationMs 5000 --offsetMs 2000
```
