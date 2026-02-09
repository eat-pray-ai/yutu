---
name: caption-insert
description: Insert (upload) a caption to a video. Use this to add new subtitles or closed captions to a video.
---

# Caption Insert

This skill provides instructions for inserting YouTube captions using the `yutu` CLI.

## Usage

```bash
yutu caption insert [flags]
```

## Options

- `--videoId`, `-v`: ID of the video.
- `--file`, `-f`: Path to the caption file to upload.
- `--language`, `-l`: Language of the caption track (BCP-47 code).
- `--name`, `-n`: Name of the caption track.
- `--audioTrackType`, `-a`: `unknown`, `primary`, `commentary`, `descriptive` (default "unknown").
- `--isAutoSynced`, `-A`: Whether YouTube synchronized the caption track (default true).
- `--isCC`, `-C`: Whether the track contains closed captions.
- `--isDraft`, `-D`: Whether the caption track is a draft.
- `--isEasyReader`, `-E`: Whether caption track is formatted for 'easy reader'.
- `--isLarge`, `-L`: Whether the caption track uses large text.
- `--trackKind`, `-t`: `standard`, `ASR`, `forced` (default "standard").
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Upload an English caption track:**

```bash
yutu caption insert --videoId VIDEO_ID --file caption.srt --language en --name "English"
```

**Upload a draft caption:**

```bash
yutu caption insert --videoId VIDEO_ID --file caption.srt --language en --isDraft true
```
