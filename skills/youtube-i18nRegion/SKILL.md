---
name: yutu-i18nRegion
description: "Manage YouTube i18n regions. Use this skill to list available internationalization regions. Always use this skill when the user mentions i18n region or wants to perform any operation on YouTube i18n region, even if they don't explicitly ask for i18n region management. Includes setup and installation instructions for first-time users. Triggers: list i18n regions, list i18n region, list my i18n region"
---

# Yutu I18n Region

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
