---
name: channelSection-list
description: List channel sections for a channel or the authenticated user.
---

# ChannelSection List

This skill provides instructions for listing YouTube channel sections using the `yutu` CLI.

## Usage

```bash
yutu channelSection list [flags]
```

## Options

- `--channelId`, `-c`: Channel ID to list sections for.
- `--mine`, `-M`: Return sections for the authenticated user.
- `--ids`, `-i`: Return sections with specified IDs.
- `--hl`, `-l`: Language code for localized content.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List my channel sections:**

```bash
yutu channelSection list --mine true
```

**List sections for another channel:**

```bash
yutu channelSection list --channelId CHANNEL_ID
```
