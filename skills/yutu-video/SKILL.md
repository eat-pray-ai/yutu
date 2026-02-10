---
name: yutu-video
description: Manage YouTube videos using the yutu CLI. Use this skill to list, upload, update, delete, rate, or report videos.
---

# Yutu Video

## Overview

This skill allows you to manage YouTube videos using the `yutu` CLI tool. You can upload videos, manage video metadata, rate videos, and delete them.

## Video Operations

### List Videos

Retrieve information about videos.

**Reference:** [references/video-list.md](references/video-list.md)

**Common Tasks:**

- List my videos: `yutu video list --myRating like`
- List most popular: `yutu video list --chart mostPopular`

### Insert/Upload Video

Upload a new video.

**Reference:** [references/video-insert.md](references/video-insert.md)

**Common Tasks:**

- Upload video: `yutu video insert --title "Title" --file video.mp4 --privacy public`

### Update Video

Update video metadata.

**Reference:** [references/video-update.md](references/video-update.md)

**Common Tasks:**

- Update title: `yutu video update --id VIDEO_ID --title "New Title"`

### Rate Video

Like or dislike a video.

**References:**

- [references/video-rate.md](references/video-rate.md)
- [references/video-getRating.md](references/video-getRating.md)

**Common Tasks:**

- Like video: `yutu video rate --ids VIDEO_ID --rating like`

### Delete Video

Remove a video.

**Reference:** [references/video-delete.md](references/video-delete.md)

**Common Tasks:**

- Delete video: `yutu video delete --ids VIDEO_ID`

### Report Abuse

Report abuse on a video.

**Reference:** [references/video-reportAbuse.md](references/video-reportAbuse.md)

## Resources

- [references/video-list.md](references/video-list.md): List videos.
- [references/video-insert.md](references/video-insert.md): Upload/insert videos.
- [references/video-update.md](references/video-update.md): Update videos.
- [references/video-rate.md](references/video-rate.md): Rate videos.
- [references/video-getRating.md](references/video-getRating.md): Get video rating.
- [references/video-reportAbuse.md](references/video-reportAbuse.md): Report abuse.
- [references/video-delete.md](references/video-delete.md): Delete videos.
