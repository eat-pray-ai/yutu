---
name: caption-update
description: Update an existing caption track. Use this to replace the caption file or update metadata like draft status.
---

# Caption Update

This skill provides instructions for updating YouTube captions using the `yutu` CLI.

## Usage

```bash
yutu caption update [flags]
```

## Options

- `--videoId`, `-v`: ID of the video.
- `--file`, `-f`: Path to the caption file.
- `--language`, `-l`: Language of the caption track.
- `--name`, `-n`: Name of the caption track.
- `--audioTrackType`, `-a`: `unknown`, `primary`, `commentary`, `descriptive`.
- `--isAutoSynced`, `-A`: Whether YouTube synchronized the caption track.
- `--isCC`, `-C`: Whether the track contains closed captions.
- `--isDraft`, `-D`: Whether the caption track is a draft.
- `--isEasyReader`, `-E`: Whether caption track is formatted for 'easy reader'.
- `--isLarge`, `-L`: Whether the caption track uses large text.
- `--trackKind`, `-t`: `standard`, `ASR`, `forced`.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update a caption file:**

```bash
yutu caption update --videoId VIDEO_ID --file new_caption.srt
```

**Publish a draft caption:**

```bash
yutu caption update --videoId VIDEO_ID --isDraft false
```
