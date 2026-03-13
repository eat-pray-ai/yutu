---
name: youtube-subscription
description: "Manage YouTube subscriptions. Use this skill to list subscriptions/subscribers, subscribe to channels, or unsubscribe. Useful when working with YouTube subscription — provides commands to delete, insert, and list subscription via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: delete subscriptions, delete subscription, delete my subscription, insert a new subscription, insert subscription, insert my subscription, list subscription information, list subscription, list my subscription"
metadata:
  openclaw:
    requires:
      env:
        - YUTU_CREDENTIAL
        - YUTU_CACHE_TOKEN
      bins:
        - yutu
      config:
        - client_secret.json
        - youtube.token.json
    primaryEnv: YUTU_CREDENTIAL
    emoji: "\U0001F3AC\U0001F430"
    homepage: https://github.com/eat-pray-ai/yutu
    install:
      - kind: node
        package: "@eat-pray-ai/yutu"
        bins: [yutu]
---

# YouTube Subscription

Manage YouTube subscriptions. Use this skill to list subscriptions/subscribers, subscribe to channels, or unsubscribe.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

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
