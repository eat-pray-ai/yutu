---
name: video-list
description: List video details.
---

# Video List

This skill provides instructions for listing YouTube video details using the `yutu` CLI.

## Usage

```bash
yutu video list [flags]
```

## Options

- `--ids`, `-i`: Return videos with the given IDs.
- `--myRating`, `-R`: Return videos liked/disliked by the authenticated user (`like`, `dislike`).
- `--chart`, `-c`: `mostPopular`.
- `--videoCategoryId`, `-g`: Filter by category ID.
- `--regionCode`, `-r`: Filter by region.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--hl`, `-l`: Localization language.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet,status]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Get video details:**

```bash
yutu video list --ids VIDEO_ID
```

**List most popular videos:**

```bash
yutu video list --chart mostPopular --regionCode US
```
