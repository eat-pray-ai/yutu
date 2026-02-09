---
name: caption-list
description: List captions of a video. Use this to see available caption tracks for a specific video.
---

# Caption List

This skill provides instructions for listing YouTube captions using the `yutu` CLI.

## Usage

```bash
yutu caption list [flags]
```

## Options

- `--videoId`, `-v`: ID of the video.
- `--ids`, `-i`: IDs of the captions to list.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--jsonpath`, `-j`: JSONPath expression to filter the output.
- `--onBehalfOf`, `-b`: Content owner ID.
- `--onBehalfOfContentOwner`, `-B`: Content owner ID.

## Examples

**List all captions for a video:**

```bash
yutu caption list --videoId VIDEO_ID
```

**List specific captions by ID:**

```bash
yutu caption list --ids "ID1,ID2"
```
