---
name: yutu-comment
description: Manage YouTube comments using the yutu CLI. Use this skill to list, insert, update, delete, mark as spam, or set moderation status for comments.
---

# Yutu Comment

## Overview

This skill allows you to manage YouTube comments using the `yutu` CLI tool. You can perform operations such as listing, replying, updating, deleting, and moderating comments.

## Comment Operations

### List Comments

Retrieve specific comments or replies.

**Reference:** [references/comment-list.md](references/comment-list.md)

**Common Tasks:**

- List replies: `yutu comment list --parentId PARENT_COMMENT_ID`

### Insert/Reply to Comments

Post new comments or replies.

**Reference:** [references/comment-insert.md](references/comment-insert.md)

**Common Tasks:**

- Post comment: `yutu comment insert --videoId VIDEO_ID --textOriginal "Hello"`
- Reply to comment: `yutu comment insert --parentId PARENT_COMMENT_ID --textOriginal "Reply"`

### Update Comments

Edit the text of an existing comment.

**Reference:** [references/comment-update.md](references/comment-update.md)

**Common Tasks:**

- Update text: `yutu comment update --id COMMENT_ID --textOriginal "New text"`

### Delete Comments

Remove a comment.

**Reference:** [references/comment-delete.md](references/comment-delete.md)

**Common Tasks:**

- Delete comment: `yutu comment delete --ids COMMENT_ID`

### Moderation

Mark comments as spam or change their moderation status (e.g., publish, reject).

**References:**

- [references/comment-markAsSpam.md](references/comment-markAsSpam.md)
- [references/comment-setModerationStatus.md](references/comment-setModerationStatus.md)

**Common Tasks:**

- Mark as spam: `yutu comment markAsSpam --ids COMMENT_ID`
- Publish comment: `yutu comment setModerationStatus --ids COMMENT_ID --moderationStatus published`

## Resources

- [references/comment-list.md](references/comment-list.md): List comments.
- [references/comment-insert.md](references/comment-insert.md): Insert comments.
- [references/comment-update.md](references/comment-update.md): Update comments.
- [references/comment-delete.md](references/comment-delete.md): Delete comments.
- [references/comment-markAsSpam.md](references/comment-markAsSpam.md): Mark comments as spam.
- [references/comment-setModerationStatus.md](references/comment-setModerationStatus.md): Set moderation status.
