---
name: video-reportAbuse
description: Report abuse on a video.
---

# Video ReportAbuse

This skill provides instructions for reporting abuse on YouTube videos using the `yutu` CLI.

## Usage

```bash
yutu video reportAbuse [flags]
```

## Options

- `--ids`, `-i`: IDs of the videos to report.
- `--reasonId`, `-r`: ID of the reason.
- `--secondaryReasonId`, `-s`: ID of the secondary reason.
- `--comments`, `-c`: Additional comments.
- `--language`, `-l`: Language content was viewed in.

## Examples

**Report a video:**

```bash
yutu video reportAbuse --ids VIDEO_ID --reasonId REASON_ID --comments "Abusive content"
```
