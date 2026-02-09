---
name: subscription-insert
description: Subscribe to a channel.
---

# Subscription Insert

This skill provides instructions for creating YouTube subscriptions using the `yutu` CLI.

## Usage

```bash
yutu subscription insert [flags]
```

## Options

- `--channelId`, `-c`: ID of the channel to subscribe to.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**Subscribe to a channel:**

```bash
yutu subscription insert --channelId CHANNEL_ID
```
