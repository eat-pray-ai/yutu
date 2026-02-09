---
name: membershipsLevel-list
description: List membership levels for the authenticated user's channel.
---

# MembershipsLevel List

This skill provides instructions for listing YouTube membership levels using the `yutu` CLI.

## Usage

```bash
yutu membershipsLevel list [flags]
```

## Options

- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List membership levels:**

```bash
yutu membershipsLevel list
```
