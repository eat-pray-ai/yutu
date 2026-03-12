---
name: youtube-video-abuse-report-reason
description: "Manage YouTube video abuse report reasons. Use this skill to list available abuse report reasons. Useful when working with YouTube video abuse report reason — covers listing, creating, updating, and deleting video abuse report reason via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list video abuse report reasons, list video abuse report reason, list my video abuse report reason"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
---

# Yutu Video Abuse Report Reason

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
