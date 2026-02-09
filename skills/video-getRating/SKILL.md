---
name: video-getRating
description: Get the rating (like/dislike/none) of a video.
---

# Video GetRating

This skill provides instructions for getting YouTube video ratings using the `yutu` CLI.

## Usage

```bash
yutu video getRating [flags]
```

## Options

- `--ids`, `-i`: IDs of the videos.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Get rating of a video:**

```bash
yutu video getRating --ids VIDEO_ID
```
