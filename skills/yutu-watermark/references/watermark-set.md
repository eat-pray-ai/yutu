# Watermark Set

Set a watermark for channel's videos. Use this skill to set a watermark for channel's videos.

## Usage

```bash
yutu watermark set [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` | Yes | ID of the channel |
| `--durationMs` | `-d` |  | Duration in milliseconds for which the watermark should be displayed |
| `--file` | `-f` | Yes | Path to the watermark file |
| `--inVideoPosition` | `-p` |  | topLeft\|topRight\|bottomLeft\|bottomRight |
| `--offsetMs` | `-m` |  | Defines the time at which the watermark will appear |
| `--offsetType` | `-t` |  | offsetFromStart\|offsetFromEnd |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Set a watermark for a channel
yutu watermark set --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file watermark.png
# Set a watermark with position and offset
yutu watermark set --channelId UC_x5X --file logo.png --inVideoPosition bottomRight --offsetType offsetFromEnd --offsetMs 1000
# Set a watermark with duration
yutu watermark set --channelId UC_x5X --file logo.png --durationMs 5000 --offsetMs 2000
```
