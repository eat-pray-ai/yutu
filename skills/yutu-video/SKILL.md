---
name: yutu-video
description: Manage YouTube videos. Use this skill when you need to list, upload, update, delete, get rating, or report videos.
---

# Yutu Video

## Overview

This skill allows you to manage YouTube videos using the `yutu` CLI tool.

## Video Operations

### Delete videos

Delete videos. Use this tool when you need to delete videos by IDs.

```bash
# Delete a video by ID
yutu video delete --ids dQw4w9WgXcQ
```

**Reference:** [references/video-delete.md](references/video-delete.md)

### Get video ratings

Get video ratings. Use this tool when you need to get video ratings by IDs.

```bash
# Get rating of a video
yutu video getRating --ids dQw4w9WgXcQ
```

**Reference:** [references/video-getRating.md](references/video-getRating.md)

### Upload a video

Upload a video. Use this tool when you need to upload a video.

```bash
# Upload a public video
yutu video insert --file video.mp4 --title 'My Video' --categoryId 22 --privacy public
```

**Reference:** [references/video-insert.md](references/video-insert.md)

### List video information

List video information. Use this tool when you need to list video information.

```bash
# List a video by ID
yutu video list --ids dQw4w9WgXcQ
```

**Reference:** [references/video-list.md](references/video-list.md)

### Rate a video

Rate a video. Use this tool when you need to rate a video.

```bash
# Like a video
yutu video rate --ids dQw4w9WgXcQ --rating like
```

**Reference:** [references/video-rate.md](references/video-rate.md)

### Report abuse on a video

Report abuse on a video. Use this tool when you need to report abuse on a video.

```bash
# Report abuse on a video
yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V
```

**Reference:** [references/video-reportAbuse.md](references/video-reportAbuse.md)

### Update a video

Update a video. Use this tool when you need to update a video.

```bash
# Update video title
yutu video update --id dQw4w9WgXcQ --title 'New Title'
```

**Reference:** [references/video-update.md](references/video-update.md)

## Resources

- [references/video-delete.md](references/video-delete.md): Detailed usage of `Delete videos`
- [references/video-getRating.md](references/video-getRating.md): Detailed usage of `Get video ratings`
- [references/video-insert.md](references/video-insert.md): Detailed usage of `Upload a video`
- [references/video-list.md](references/video-list.md): Detailed usage of `List video information`
- [references/video-rate.md](references/video-rate.md): Detailed usage of `Rate a video`
- [references/video-reportAbuse.md](references/video-reportAbuse.md): Detailed usage of `Report abuse on a video`
- [references/video-update.md](references/video-update.md): Detailed usage of `Update a video`
