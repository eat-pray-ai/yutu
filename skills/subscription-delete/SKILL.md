---
name: subscription-delete
description: Delete a subscription (unsubscribe).
---

# Subscription Delete

This skill provides instructions for deleting YouTube subscriptions using the `yutu` CLI.

## Usage

```bash
yutu subscription delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the subscriptions to delete.

## Examples

**Unsubscribe from a channel:**

```bash
yutu subscription delete --ids SUBSCRIPTION_ID
```
