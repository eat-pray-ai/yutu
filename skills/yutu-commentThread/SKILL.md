---
name: yutu-commentThread
description: Manage YouTube comment threads. Use this skill when you need to list or insert new top-level comment threads.
---

# Yutu CommentThread

## Overview

This skill allows you to manage YouTube comment threads using the `yutu` CLI tool.

## CommentThread Operations

### Insert a new comment thread

Insert a new comment thread. Use this tool when you need to insert a new comment thread.

```bash
# Post a comment on a video
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Great video!'
```

**Reference:** [references/commentThread-insert.md](references/commentThread-insert.md)

### List comment threads

List comment threads. Use this tool when you need to list comment threads.

```bash
# List comment threads on a video
yutu commentThread list --videoId dQw4w9WgXcQ --maxResults 10
```

**Reference:** [references/commentThread-list.md](references/commentThread-list.md)

## Resources

- [references/commentThread-insert.md](references/commentThread-insert.md): Detailed usage of `Insert a new comment thread`
- [references/commentThread-list.md](references/commentThread-list.md): Detailed usage of `List comment threads`
