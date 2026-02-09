---
name: superChatEvent-list
description: List Super Chat events for a channel.
---

# SuperChatEvent List

This skill provides instructions for listing YouTube Super Chat events using the `yutu` CLI.

## Usage

```bash
yutu superChatEvent list [flags]
```

## Options

- `--hl`, `-l`: Return rendered funding amounts in specified language.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List Super Chat events:**

```bash
yutu superChatEvent list
```
