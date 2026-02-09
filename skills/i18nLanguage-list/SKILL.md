---
name: i18nLanguage-list
description: List supported i18n languages.
---

# I18nLanguage List

This skill provides instructions for listing YouTube i18n languages using the `yutu` CLI.

## Usage

```bash
yutu i18nLanguage list [flags]
```

## Options

- `--hl`, `-l`: Host language.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List supported languages for en_US:**

```bash
yutu i18nLanguage list --hl en_US
```
