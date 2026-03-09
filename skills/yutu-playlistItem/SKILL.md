---
name: yutu-playlistItem
description: "Manage YouTube playlist items. Use this skill to list items in a playlist, add new items, update items, or remove items. Always use this skill when the user mentions playlist item or wants to perform any operation on YouTube playlist item, even if they don't explicitly ask for playlist item management. Triggers: delete items from a playlist, delete playlist item, delete my playlist item, insert a playlist item into a playlist, insert playlist item, insert my playlist item, list playlist items, list playlist item, list my playlist item, update a playlist item, update playlist item, update my playlist item"
---

# Yutu Playlist Item

Manage YouTube playlist items. Use this skill to list items in a playlist, add new items, update items, or remove items.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete items from a playlist | [details](references/playlistItem-delete.md) |
| insert | Insert a playlist item into a playlist | [details](references/playlistItem-insert.md) |
| list | List playlist items | [details](references/playlistItem-list.md) |
| update | Update a playlist item | [details](references/playlistItem-update.md) |

## Quick Start

```bash
# Show all playlist item commands
yutu playlistItem --help

# List playlist item
yutu playlistItem list
```
