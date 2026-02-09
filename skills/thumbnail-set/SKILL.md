---
name: thumbnail-set
description: Set a custom thumbnail for a video.
---

# Thumbnail Set

This skill provides instructions for setting YouTube video thumbnails using the `yutu` CLI.

## Usage

```bash
yutu thumbnail set [flags]
```

## Options

- `--videoId`, `-v`: ID of the video.
- `--file`, `-f`: Path to the thumbnail file.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Set video thumbnail:**

```bash
yutu thumbnail set --videoId VIDEO_ID --file thumbnail.jpg
```
