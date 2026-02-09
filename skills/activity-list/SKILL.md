---
name: activity-list
description: List YouTube activities (likes, favorites, uploads, etc.) for a channel or the authenticated user. Use this when the user asks to see recent activities, check uploads, or monitor channel interactions.
---

# Activity List

This skill provides instructions for listing YouTube activities using the `yutu` CLI directly.

## Usage

Execute the command from the repository root:

```bash
yutu activity list [flags]
```

## Options

- `--channelId`, `-c`: ID of the channel.
- `--home`, `-H`: List home activities (default: true).
- `--mine`, `-M`: List activities for the authenticated user (default: true).
- `--maxResults`, `-n`: Maximum number of items to return (default: 5). Set to 0 for no limit.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma-separated parts to retrieve (default: `[id,snippet,contentDetails]`).
- `--publishedAfter`, `-a`: Filter activities published after this date.
- `--publishedBefore`, `-b`: Filter activities published before this date.
- `--regionCode`, `-r`: Region code.
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List my recent activities (default):**

```bash
yutu activity list --mine true
```

**List activities for a specific channel:**

```bash
yutu activity list --channelId CHANNEL_ID --mine false
```

**List 10 activities in JSON format:**

```bash
yutu activity list --maxResults 10 --output json
```
