---
name: comment-list
description: List comments by IDs or list replies to a parent comment.
---

# Comment List

This skill provides instructions for listing YouTube comments using the `yutu` CLI.

## Usage

```bash
yutu comment list [flags]
```

## Options

- `--parentId`, `-P`: ID of the parent comment (to list replies).
- `--ids`, `-i`: IDs of comments to list.
- `--textFormat`, `-t`: `textFormatUnspecified`, `html`, `plainText` (default "html").
- `--maxResults`, `-n`: Maximum number of items to return.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List replies to a comment:**

```bash
yutu comment list --parentId PARENT_COMMENT_ID
```

**Get a comment by ID:**

```bash
yutu comment list --ids COMMENT_ID
```
