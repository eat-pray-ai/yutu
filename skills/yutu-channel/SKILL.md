---
name: yutu-channel
description: Manage YouTube channels using the yutu CLI. Use this skill when you need to list channel information (search by handle, username, id, or mine) or update channel details (title, description, country, etc.).
---

# Yutu Channel

## Overview

This skill allows you to manage YouTube channels using the `yutu` CLI tool. You can retrieve detailed information about channels and update channel metadata.

## Channel Operations

### List Channels

Retrieve information about channels. You can find channels by ID, handle, username, or list your own channels.

**Reference:** [references/channel-list.md](references/channel-list.md)

**Common Tasks:**

- List my channel: `yutu channel list --mine`
- Find channel by handle: `yutu channel list --forHandle @username`
- Find channel by ID: `yutu channel list --ids CHANNEL_ID`

### Update Channels

Update channel metadata such as title, description, and default language.

**Reference:** [references/channel-update.md](references/channel-update.md)

**Common Tasks:**

- Update title: `yutu channel update --id CHANNEL_ID --title "New Title"`
- Update description: `yutu channel update --id CHANNEL_ID --description "New Description"`

## Resources

- [references/channel-list.md](references/channel-list.md): Detailed flags and usage for listing channels.
- [references/channel-update.md](references/channel-update.md): Detailed flags and usage for updating channels.
