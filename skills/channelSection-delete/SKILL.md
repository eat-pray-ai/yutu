---
name: channelSection-delete
description: Delete channel sections by IDs. Use this to reorganize a channel page.
---

# ChannelSection Delete

This skill provides instructions for deleting YouTube channel sections using the `yutu` CLI.

## Usage

```bash
yutu channelSection delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the channel sections to delete (comma-separated).
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.

## Examples

**Delete a channel section:**

```bash
yutu channelSection delete --ids SECTION_ID
```
