---
name: comment-markAsSpam
description: Mark comments as spam by IDs.
---

# Comment MarkAsSpam

This skill provides instructions for marking YouTube comments as spam using the `yutu` CLI.

## Usage

```bash
yutu comment markAsSpam [flags]
```

## Options

- `--ids`, `-i`: IDs of the comments to mark as spam.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Mark a comment as spam:**

```bash
yutu comment markAsSpam --ids COMMENT_ID
```
