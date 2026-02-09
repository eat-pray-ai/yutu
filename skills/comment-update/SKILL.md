---
name: comment-update
description: Update an existing comment (e.g., change text).
---

# Comment Update

This skill provides instructions for updating YouTube comments using the `yutu` CLI.

## Usage

```bash
yutu comment update [flags]
```

## Options

- `--id`, `-i`: ID of the comment to update.
- `--textOriginal`, `-t`: New text of the comment.
- `--viewerRating`, `-r`: `none`, `like`, `dislike`.
- `--canRate`, `-R`: Whether the viewer can rate the comment.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update comment text:**

```bash
yutu comment update --id COMMENT_ID --textOriginal "Updated text."
```
