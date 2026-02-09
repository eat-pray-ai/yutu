---
name: channel-list
description: List channel info (title, description, stats) for specific channels or the authenticated user's channel.
---

# Channel List

This skill provides instructions for listing YouTube channels using the `yutu` CLI.

## Usage

```bash
yutu channel list [flags]
```

## Options

- `--mine`, `-M`: Return the channels owned by the authenticated user (default true).
- `--ids`, `-i`: Return the channels with the specified IDs.
- `--forHandle`, `-d`: Return the channel associated with a YouTube handle.
- `--forUsername`, `-u`: Return the channel associated with a YouTube username.
- `--managedByMe`, `-E`: Return the channels managed by the authenticated user.
- `--mySubscribers`, `-S`: Return the channels subscribed to the authenticated user.
- `--categoryId`, `-g`: Return the channels within the specified guide category ID.
- `--hl`, `-l`: Localization language of the metadata.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet,status]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List my channel info:**

```bash
yutu channel list --mine true
```

**Get channel info by ID:**

```bash
yutu channel list --ids CHANNEL_ID --mine false
```

**Get channel info by handle:**

```bash
yutu channel list --forHandle "@Handle" --mine false
```
