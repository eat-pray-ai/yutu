---
name: youtube-i18n-language
description: "Manage YouTube i18n languages. Use this skill to list available internationalization languages. Useful when working with YouTube i18n language — provides commands to list i18n language via the yutu CLI. Includes setup and installation instructions for first-time users. Triggers: list i18n languages, list i18n language, list my i18n language"
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

# YouTube I18n Language

Manage YouTube i18n languages. Use this skill to list available internationalization languages.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Read the linked reference for full flag details and examples.

| Operation | Description | Reference |
|-----------|-------------|----------|
| list | List i18n languages | [details](references/i18nLanguage-list.md) |

## Quick Start

```bash
# Show all i18n language commands
yutu i18nLanguage --help

# List i18n language
yutu i18nLanguage list
```
