---
name: watermark-set
description: Set a channel watermark.
---

# Watermark Set

This skill provides instructions for setting a YouTube channel watermark using the `yutu` CLI.

## Usage

```bash
yutu watermark set [flags]
```

## Options

- `--channelId`, `-c`: ID of the channel.
- `--file`, `-f`: Path to the watermark file.
- `--inVideoPosition`, `-p`: `topLeft`, `topRight`, `bottomLeft`, `bottomRight`.
- `--offsetType`, `-t`: `offsetFromStart`, `offsetFromEnd`.
- `--offsetMs`, `-m`: Offset in milliseconds.
- `--durationMs`, `-d`: Duration in milliseconds.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.

## Examples

**Set channel watermark:**

```bash
yutu watermark set --channelId CHANNEL_ID --file watermark.png --inVideoPosition bottomRight
```
