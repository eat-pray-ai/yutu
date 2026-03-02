---
name: yutu-subscription
description: Manage YouTube subscriptions. Use this skill when you need to list subscriptions/subscribers, subscribe to channels, or unsubscribe.
---

# Yutu Subscription

## Overview

This skill allows you to manage YouTube subscriptions using the `yutu` CLI tool.

## Subscription Operations

### Delete subscriptions

Delete subscriptions. Use this tool when you need to delete subscriptions by IDs.

```bash
# Delete a subscription by ID
yutu subscription delete --ids abc123
```

**Reference:** [references/subscription-delete.md](references/subscription-delete.md)

### Insert a new subscription

Insert a new subscription. Use this tool when you need to insert a new subscription.

```bash
# Subscribe to a channel
yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
```

**Reference:** [references/subscription-insert.md](references/subscription-insert.md)

### List subscription information

List subscription information. Use this tool when you need to list subscription information.

```bash
# List my subscriptions
yutu subscription list --mine
```

**Reference:** [references/subscription-list.md](references/subscription-list.md)

## Resources

- [references/subscription-delete.md](references/subscription-delete.md): Detailed usage of `Delete subscriptions`
- [references/subscription-insert.md](references/subscription-insert.md): Detailed usage of `Insert a new subscription`
- [references/subscription-list.md](references/subscription-list.md): Detailed usage of `List subscription information`
