---
name: caption-download
description: Download caption from a video. Use this to save caption tracks to local files (srt, vtt, sbv).
---

# Caption Download

This skill provides instructions for downloading YouTube captions using the `yutu` CLI.

## Usage

```bash
yutu caption download [flags]
```

## Options

- `--file`, `-f`: Path to save the caption file.
- `--id`, `-i`: ID of the caption to download.
- `--onBehalfOf`, `-b`: Content owner ID.
- `--onBehalfOfContentOwner`, `-B`: Content owner ID.
- `--tfmt`, `-t`: Caption format: `sbv`, `srt`, `vtt`.
- `--tlang`, `-l`: Translate the captions into this language.

## Examples

**Download a caption as SRT:**

```bash
yutu caption download --id CAPTION_ID --file caption.srt --tfmt srt
```

**Download and translate to Spanish:**

```bash
yutu caption download --id CAPTION_ID --tlang es --tfmt srt --file caption_es.srt
```
