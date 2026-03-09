---
name: yutu-subscription
description: "Manage YouTube subscriptions. Use this skill to list subscriptions/subscribers, subscribe to channels, or unsubscribe. Always use this skill when the user mentions subscription or wants to perform any operation on YouTube subscription, even if they don't explicitly ask for subscription management. Triggers: delete subscriptions, delete subscription, delete my subscription, insert a new subscription, insert subscription, insert my subscription, list subscription information, list subscription, list my subscription"
---

# Yutu Subscription

Manage YouTube subscriptions. Use this skill to list subscriptions/subscribers, subscribe to channels, or unsubscribe.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| delete | Delete subscriptions | [details](references/subscription-delete.md) |
| insert | Insert a new subscription | [details](references/subscription-insert.md) |
| list | List subscription information | [details](references/subscription-list.md) |

## Quick Start

```bash
# Show all subscription commands
yutu subscription --help

# List subscription
yutu subscription list
```
