---
name: subscription-list
description: List subscriptions for a channel or the authenticated user.
---

# Subscription List

This skill provides instructions for listing YouTube subscriptions using the `yutu` CLI.

## Usage

```bash
yutu subscription list [flags]
```

## Options

- `--mine`, `-M`: Return subscriptions of the authenticated user (default true).
- `--channelId`, `-c`: Return subscriptions of the given channel owner.
- `--mySubscribers`, `-S`: Return the subscribers of the authenticated user.
- `--ids`, `-i`: Return the subscriptions with the given IDs.
- `--forChannelId`, `-C`: Return subscriptions to the subset of these channels.
- `--maxResults`, `-n`: Maximum number of items to return.
- `--order`, `-O`: `relevance`, `unread`, `alphabetical`.
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--parts`, `-p`: Comma separated parts (default `[id,snippet]`).
- `--jsonpath`, `-j`: JSONPath expression to filter the output.

## Examples

**List my subscriptions:**

```bash
yutu subscription list --mine true
```

**List my subscribers:**

```bash
yutu subscription list --mySubscribers true
```
