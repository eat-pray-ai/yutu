---
name: comment-insert
description: Reply to a comment or add a top-level comment (note: use commentThread insert for top-level usually, but check API). Actually, `comment insert` is used for replying to a thread (parentId) or specific context.
---

# Comment Insert

This skill provides instructions for inserting a YouTube comment using the `yutu` CLI.
Typically used for replying to an existing comment (parent).

## Usage

```bash
yutu comment insert [flags]
```

## Options

- `--parentId`, `-P`: ID of the parent comment (required for replies).
- `--textOriginal`, `-t`: Text of the comment.
- `--videoId`, `-v`: ID of the video (optional if parentId implies context).
- `--channelId`, `-c`: Channel ID of the video owner.
- `--authorChannelId`, `-a`: Channel ID of the comment author.
- `--canRate`, `-R`: Whether the viewer can rate the comment.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Reply to a comment:**

```bash
yutu comment insert --parentId PARENT_COMMENT_ID --textOriginal "This is a reply."
```
