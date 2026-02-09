---
name: videoAbuseReportReason-list
description: List abuse report reasons.
---

# VideoAbuseReportReason List

This skill provides instructions for listing YouTube video abuse report reasons using the `yutu` CLI.

## Usage

```bash
yutu videoAbuseReportReason list [flags]
```

## Options

- `--hl`, `-l`: Host language.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List abuse report reasons:**

```bash
yutu videoAbuseReportReason list --hl en_US
```
