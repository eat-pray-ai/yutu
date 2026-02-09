---
name: video-rate
description: Rate a video (like/dislike/none).
---

# Video Rate

This skill provides instructions for rating YouTube videos using the `yutu` CLI.

## Usage

```bash
yutu video rate [flags]
```

## Options

- `--ids`, `-i`: IDs of the videos to rate.
- `--rating`, `-r`: `like`, `dislike`, `none`.

## Examples

**Like a video:**

```bash
yutu video rate --ids VIDEO_ID --rating like
```
