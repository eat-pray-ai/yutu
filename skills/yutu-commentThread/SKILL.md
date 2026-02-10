---
name: yutu-commentThread
description: Manage YouTube comment threads using the yutu CLI. Use this skill to list or insert new top-level comment threads.
---

# Yutu Comment Thread

## Overview

This skill allows you to manage YouTube comment threads using the `yutu` CLI tool. Comment threads are top-level comments and their replies.

## Comment Thread Operations

### List Comment Threads

Retrieve a list of comment threads for a video or channel.

**Reference:** [references/commentThread-list.md](references/commentThread-list.md)

**Common Tasks:**

- List threads for video: `yutu commentThread list --videoId VIDEO_ID`

### Insert Comment Thread

Create a new top-level comment thread on a video or channel.

**Reference:** [references/commentThread-insert.md](references/commentThread-insert.md)

**Common Tasks:**

- New comment thread: `yutu commentThread insert --videoId VIDEO_ID --textOriginal "This is a new thread"`

## Resources

- [references/commentThread-list.md](references/commentThread-list.md): List comment threads.
- [references/commentThread-insert.md](references/commentThread-insert.md): Insert comment threads.
