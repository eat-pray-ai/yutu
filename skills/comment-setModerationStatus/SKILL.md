---
name: comment-setModerationStatus
description: Set moderation status (heldForReview, published, rejected) for comments.
---

# Comment SetModerationStatus

This skill provides instructions for setting the moderation status of YouTube comments using the `yutu` CLI.

## Usage

```bash
yutu comment setModerationStatus [flags]
```

## Options

- `--ids`, `-i`: IDs of the comments.
- `--moderationStatus`, `-s`: Status: `heldForReview`, `published`, `rejected`.
- `--banAuthor`, `-A`: If set to true, the author of the comment gets added to the ban list.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Publish a held comment:**

```bash
yutu comment setModerationStatus --ids COMMENT_ID --moderationStatus published
```

**Reject a comment and ban author:**

```bash
yutu comment setModerationStatus --ids COMMENT_ID --moderationStatus rejected --banAuthor true
```
