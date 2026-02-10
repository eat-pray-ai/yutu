---
name: yutu-subscription
description: Manage YouTube subscriptions using the yutu CLI. Use this skill to list subscriptions/subscribers, subscribe to channels, or unsubscribe.
---

# Yutu Subscription

## Overview

This skill allows you to manage YouTube subscriptions using the `yutu` CLI tool. You can list your subscriptions, see who subscribed to you, subscribe to new channels, and unsubscribe.

## Subscription Operations

### List Subscriptions

Retrieve information about subscriptions.

**Reference:** [references/subscription-list.md](references/subscription-list.md)

**Common Tasks:**

- List my subscriptions: `yutu subscription list --mine`
- List my subscribers: `yutu subscription list --mySubscribers`

### Insert/Subscribe

Subscribe to a channel.

**Reference:** [references/subscription-insert.md](references/subscription-insert.md)

**Common Tasks:**

- Subscribe: `yutu subscription insert --channelId CHANNEL_ID`

### Delete/Unsubscribe

Remove a subscription.

**Reference:** [references/subscription-delete.md](references/subscription-delete.md)

**Common Tasks:**

- Unsubscribe: `yutu subscription delete --ids SUBSCRIPTION_ID`

## Resources

- [references/subscription-list.md](references/subscription-list.md): List subscriptions.
- [references/subscription-insert.md](references/subscription-insert.md): Insert subscriptions.
- [references/subscription-delete.md](references/subscription-delete.md): Delete subscriptions.
