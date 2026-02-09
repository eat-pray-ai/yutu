---
name: videoCategory-list
description: List video categories.
---

# VideoCategory List

This skill provides instructions for listing YouTube video categories using the `yutu` CLI.

## Usage

```bash
yutu videoCategory list [flags]
```

## Options

- `--regionCode`, `-r`: Region code (default "US").
- `--ids`, `-i`: IDs of the video categories.
- `--hl`, `-l`: Host language.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List video categories for US:**

```bash
yutu videoCategory list --regionCode US
```
