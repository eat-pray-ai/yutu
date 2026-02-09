---
name: commentThread-list
description: List comment threads for a video or channel.
---

# CommentThread List

This skill provides instructions for listing YouTube comment threads using the `yutu` CLI.

## Usage

```bash
yutu commentThread list [flags]
```

## Options

- `--videoId`, `-v`: Returns the comment threads of the specified video.
- `--channelId`, `-c`: Returns the comment threads of the specified channel.
- `--allThreadsRelatedToChannelId`, `-a`: Returns threads of all videos of the channel + channel comments.
- `--ids`, `-i`: IDs of the comment threads.
- `--searchTerms`, `-s`: Search terms.
- `--moderationStatus`, `-m`: `published`, `heldForReview`, `likelySpam`, `rejected` (default "published").
- `--order`, `-O`: `time`, `relevance` (default "time").
- `--textFormat`, `-t`: `textFormatUnspecified`, `html` (default "html").
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List comments for a video:**

```bash
yutu commentThread list --videoId VIDEO_ID
```

**Search comments on a video:**

```bash
yutu commentThread list --videoId VIDEO_ID --searchTerms "awesome"
```
