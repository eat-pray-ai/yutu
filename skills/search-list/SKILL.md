---
name: search-list
description: Search for videos, channels, or playlists on YouTube.
---

# Search List

This skill provides instructions for searching YouTube resources using the `yutu` CLI.

## Usage

```bash
yutu search list [flags]
```

## Options

- `--q`: Textual search terms to match.
- `--maxResults`, `-n`: Maximum number of items to return (default 5).
- `--types`: Restrict results to a particular set of resource types (video, channel, playlist).
- `--channelId`: Filter on resources belonging to this channelId.
- `--publishedAfter`: Filter on resources published after this date.
- `--publishedBefore`: Filter on resources published before this date.
- `--order`: `date`, `rating`, `viewCount`, `relevance`, `title`, `videoCount` (default "relevance").
- `--output`, `-o`: Output format (`json`, `yaml`, `table`). Default: `table`.
- `--jsonpath`: JSONPath expression to filter the output.

## Examples

**Search for "golang" videos:**

```bash
yutu search list --q "golang" --types video
```

**Search for channels by name:**

```bash
yutu search list --q "Google Developers" --types channel
```
