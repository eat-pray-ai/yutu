---
name: youtube-video-abuse-report-reason
description: "Manage YouTube video abuse report reasons. Use this skill to list available abuse report reasons. Useful when working with YouTube video abuse report reason — provides commands to list video abuse report reason via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list video abuse report reasons, list video abuse report reason, list my video abuse report reason"
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

# YouTube Video Abuse Report Reason

Manage YouTube video abuse report reasons. Use this skill to list available abuse report reasons.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List video abuse report reasons | [details](references/videoAbuseReportReason-list.md) |

## Quick Start

```bash
# Show all video abuse report reason commands
yutu videoAbuseReportReason --help

# List video abuse report reason
yutu videoAbuseReportReason list
```
