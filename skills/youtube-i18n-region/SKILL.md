---
name: youtube-i18n-region
description: "Manage YouTube i18n regions. Use this skill to list available internationalization regions. Useful when working with YouTube i18n region — provides commands to list i18n region via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list i18n regions, list i18n region, list my i18n region"
compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.
metadata:
  author: eat-pray-ai
  required_config_paths:
    - client_secret.json
    - youtube.token.json
  env:
    - YUTU_CREDENTIAL
    - YUTU_CACHE_TOKEN
---

# YouTube I18n Region

Manage YouTube i18n regions. Use this skill to list available internationalization regions.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List i18n regions | [details](references/i18nRegion-list.md) |

## Quick Start

```bash
# Show all i18n region commands
yutu i18nRegion --help

# List i18n region
yutu i18nRegion list
```
