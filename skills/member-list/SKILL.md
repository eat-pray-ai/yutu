---
name: member-list
description: List members of a channel (requires authenticated user to be the channel owner).
---

# Member List

This skill provides instructions for listing YouTube channel members using the `yutu` CLI.

## Usage

```bash
yutu member list [flags]
```

## Options

- `--mode`, `-m`: `listMembersModeUnknown`, `updates`, `all_current` (default "all_current").
- `--memberChannelId`, `-c`: Comma separated list of member channel IDs to filter.
- `--hasAccessToLevel`, `-a`: Filter members who have access to a specific level.
- `--maxResults`, `-n`: Maximum number of items to return.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List current members:**

```bash
yutu member list --mode all_current
```
