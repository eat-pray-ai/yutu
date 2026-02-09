---
name: commentThread-insert
description: Insert a new comment thread. Use this to post a top-level comment on a video or channel.
---

# CommentThread Insert

This skill provides instructions for inserting a YouTube comment thread using the `yutu` CLI.

## Usage

```bash
yutu commentThread insert [flags]
```

## Options

- `--videoId`, `-v`: ID of the video to comment on.
- `--channelId`, `-c`: Channel ID (to post on channel discussion).
- `--textOriginal`, `-t`: Text of the comment.
- `--authorChannelId`, `-a`: Channel ID of the comment author.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Post a comment on a video:**

```bash
yutu commentThread insert --videoId VIDEO_ID --textOriginal "Great video!"
```

**Post a comment on a channel:**

```bash
yutu commentThread insert --channelId CHANNEL_ID --textOriginal "Hello channel!"
```
