---
name: yutu-comment
description: Manage YouTube comments. Use this skill when you need to list, create, update, delete, mark as spam, or set moderation status for comments.
---

# Yutu Comment

## Overview

This skill allows you to manage YouTube comments using the `yutu` CLI tool.

## Comment Operations

### Delete comments

Delete comments. Use this tool when you need to delete comments by IDs.

```bash
# Delete a comment by ID
yutu comment delete --ids abc123
```

**Reference:** [references/comment-delete.md](references/comment-delete.md)

### Create a comment

Create a comment. Use this tool when you need to create a comment on a video.

```bash
# Reply to a comment
yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Hello'
```

**Reference:** [references/comment-insert.md](references/comment-insert.md)

### List comments

List comments. Use this tool when you need to list comments by IDs.

```bash
# List replies to a comment
yutu comment list --parentId UgyXXXXXXXX --maxResults 10
```

**Reference:** [references/comment-list.md](references/comment-list.md)

### Mark comments as spam

Mark comments as spam. Use this tool when you need to mark comments as spam.

```bash
# Mark a comment as spam
yutu comment markAsSpam --ids abc123
```

**Reference:** [references/comment-markAsSpam.md](references/comment-markAsSpam.md)

### Set comment moderation status

Set comment moderation status. Use this tool when you need to set comment moderation status.

```bash
# Publish a held comment
yutu comment setModerationStatus --ids abc123 --moderationStatus published
```

**Reference:** [references/comment-setModerationStatus.md](references/comment-setModerationStatus.md)

### Update a comment on a video

Update a comment on a video. Use this tool when you need to update a comment on a video.

```bash
# Update comment text
yutu comment update --id abc123 --textOriginal 'Updated comment'
```

**Reference:** [references/comment-update.md](references/comment-update.md)

## Resources

- [references/comment-delete.md](references/comment-delete.md): Detailed usage of `Delete comments`
- [references/comment-insert.md](references/comment-insert.md): Detailed usage of `Create a comment`
- [references/comment-list.md](references/comment-list.md): Detailed usage of `List comments`
- [references/comment-markAsSpam.md](references/comment-markAsSpam.md): Detailed usage of `Mark comments as spam`
- [references/comment-setModerationStatus.md](references/comment-setModerationStatus.md): Detailed usage of `Set comment moderation status`
- [references/comment-update.md](references/comment-update.md): Detailed usage of `Update a comment on a video`
