---
name: yutu-channelSection
description: Manage YouTube channel sections. Use this skill when you need to list or delete channel sections.
---

# Yutu ChannelSection

## Overview

This skill allows you to manage YouTube channel sections using the `yutu` CLI tool.

## ChannelSection Operations

### Delete channel sections

Delete channel sections. Use this tool when you need to delete channel sections by IDs.

```bash
# Delete a channel section by ID
yutu channelSection delete --ids abc123
```

**Reference:** [references/channelSection-delete.md](references/channelSection-delete.md)

### List channel sections

List channel sections. Use this tool when you need to list channel sections.

```bash
# List my channel sections
yutu channelSection list --mine
```

**Reference:** [references/channelSection-list.md](references/channelSection-list.md)

## Resources

- [references/channelSection-delete.md](references/channelSection-delete.md): Detailed usage of `Delete channel sections`
- [references/channelSection-list.md](references/channelSection-list.md): Detailed usage of `List channel sections`
