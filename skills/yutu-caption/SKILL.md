---
name: yutu-caption
description: Manage YouTube video captions. Use this skill when you need to list, insert, update, download, or delete video captions.
---

# Yutu Caption

## Overview

This skill allows you to manage YouTube captions using the `yutu` CLI tool.

## Caption Operations

### Delete captions

Delete captions. Use this tool when you need to delete captions of a video by IDs.

```bash
# Delete a caption by ID
yutu caption delete --ids abc123
```

**Reference:** [references/caption-delete.md](references/caption-delete.md)

### Download a caption

Download a caption. Use this tool when you need to download a caption from a video.

```bash
# Download a caption as SRT
yutu caption download --id abc123 --file subtitle.srt
```

**Reference:** [references/caption-download.md](references/caption-download.md)

### Insert a caption

Insert a caption. Use this tool when you need to insert a caption to a video.

```bash
# Insert a caption to a video
yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ
```

**Reference:** [references/caption-insert.md](references/caption-insert.md)

### List captions

List captions. Use this tool when you need to list captions of a video.

```bash
# List captions of a video
yutu caption list --videoId dQw4w9WgXcQ
```

**Reference:** [references/caption-list.md](references/caption-list.md)

### Update a video caption

Update a video caption. Use this tool when you need to update a video caption.

```bash
# Publish a draft caption
yutu caption update --videoId dQw4w9WgXcQ --isDraft=false
```

**Reference:** [references/caption-update.md](references/caption-update.md)

## Resources

- [references/caption-delete.md](references/caption-delete.md): Detailed usage of `Delete captions`
- [references/caption-download.md](references/caption-download.md): Detailed usage of `Download a caption`
- [references/caption-insert.md](references/caption-insert.md): Detailed usage of `Insert a caption`
- [references/caption-list.md](references/caption-list.md): Detailed usage of `List captions`
- [references/caption-update.md](references/caption-update.md): Detailed usage of `Update a video caption`
