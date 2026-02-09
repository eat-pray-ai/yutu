---
name: channel-update
description: Update channel info (title, description, country, custom URL).
---

# Channel Update

This skill provides instructions for updating YouTube channels using the `yutu` CLI.

## Usage

```bash
yutu channel update [flags]
```

## Options

- `--id`, `-i`: ID of the channel to update.
- `--title`, `-t`: Title of the channel.
- `--description`, `-d`: Description of the channel.
- `--country`, `-c`: Country of the channel.
- `--customUrl`, `-u`: Custom URL of the channel.
- `--defaultLanguage`, `-l`: The language of the channel's default title and description.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update channel title and description:**

```bash
yutu channel update --id CHANNEL_ID --title "New Title" --description "New Description"
```
