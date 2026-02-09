---
name: channelBanner-insert
description: Upload and insert a new banner image for a channel.
---

# ChannelBanner Insert

This skill provides instructions for inserting a YouTube channel banner using the `yutu` CLI.

## Usage

```bash
yutu channelBanner insert [flags]
```

## Options

- `--channelId`, `-c`: ID of the channel to insert the banner for.
- `--file`, `-f`: Path to the banner image file.
- `--onBehalfOfContentOwner`, `-b`: Content owner ID.
- `--onBehalfOfContentOwnerChannel`, `-B`: Content owner channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Upload a channel banner:**

```bash
yutu channelBanner insert --channelId CHANNEL_ID --file banner.png
```
