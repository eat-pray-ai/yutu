---
name: caption-delete
description: Delete captions of a video by ids. Use this when the user wants to remove specific caption tracks.
---

# Caption Delete

This skill provides instructions for deleting YouTube captions using the `yutu` CLI.

## Usage

```bash
yutu caption delete [flags]
```

## Options

- `--ids`, `-i`: IDs of the captions to delete (comma-separated strings).
- `--onBehalfOf`, `-b`: Content owner ID.
- `--onBehalfOfContentOwner`, `-B`: Content owner ID.

## Examples

**Delete a single caption:**

```bash
yutu caption delete --ids CAPTION_ID
```

**Delete multiple captions:**

```bash
yutu caption delete --ids "ID1,ID2"
```
