---
name: yutu-caption
description: Manage YouTube captions using the yutu CLI. Use this skill to list, insert, update, download, or delete video captions.
---

# Yutu Caption

## Overview

This skill allows you to manage YouTube video captions using the `yutu` CLI tool. You can perform operations such as listing, uploading (inserting), updating, downloading, and deleting captions.

## Caption Operations

### List Captions

Retrieve a list of captions available for a specific video.

**Reference:** [references/caption-list.md](references/caption-list.md)

**Common Tasks:**

- List captions: `yutu caption list --videoId VIDEO_ID`

### Download Captions

Download specific caption tracks.

**Reference:** [references/caption-download.md](references/caption-download.md)

**Common Tasks:**

- Download caption: `yutu caption download --id CAPTION_ID --tfmt srt --file output.srt`

### Insert/Upload Captions

Upload new caption files to a video.

**Reference:** [references/caption-insert.md](references/caption-insert.md)

**Common Tasks:**

- Upload caption: `yutu caption insert --videoId VIDEO_ID --language en --name "English" --file captions.srt`

### Update Captions

Update existing caption tracks (e.g., replace file or update metadata).

**Reference:** [references/caption-update.md](references/caption-update.md)

### Delete Captions

Remove caption tracks from a video.

**Reference:** [references/caption-delete.md](references/caption-delete.md)

**Common Tasks:**

- Delete caption: `yutu caption delete --ids CAPTION_ID`

## Resources

- [references/caption-list.md](references/caption-list.md): List captions.
- [references/caption-download.md](references/caption-download.md): Download captions.
- [references/caption-insert.md](references/caption-insert.md): Insert captions.
- [references/caption-update.md](references/caption-update.md): Update captions.
- [references/caption-delete.md](references/caption-delete.md): Delete captions.
