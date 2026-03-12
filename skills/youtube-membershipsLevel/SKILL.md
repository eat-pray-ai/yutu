---
name: yutu-membershipsLevel
description: "Manage YouTube memberships levels. Use this skill to list information about channel membership levels. Always use this skill when the user mentions memberships level or wants to perform any operation on YouTube memberships level, even if they don't explicitly ask for memberships level management. Includes setup and installation instructions for first-time users. Triggers: list memberships levels, list memberships level, list my memberships level"
---

# Yutu Memberships Level

Manage YouTube memberships levels. Use this skill to list information about channel membership levels.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List memberships levels | [details](references/membershipsLevel-list.md) |

## Quick Start

```bash
# Show all memberships level commands
yutu membershipsLevel --help

# List memberships level
yutu membershipsLevel list
```
