---
name: i18nRegion-list
description: List supported i18n regions.
---

# I18nRegion List

This skill provides instructions for listing YouTube i18n regions using the `yutu` CLI.

## Usage

```bash
yutu i18nRegion list [flags]
```

## Options

- `--hl`, `-l`: Host language.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List supported regions for en_US:**

```bash
yutu i18nRegion list --hl en_US
```
