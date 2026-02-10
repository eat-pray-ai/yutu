---
name: yutu-watermark
description: Manage YouTube watermarks using the yutu CLI. Use this skill to set or unset watermarks for channel videos.
---

# Yutu Watermark

## Overview

This skill allows you to manage YouTube watermarks using the `yutu` CLI tool. You can set a branding watermark for all your videos or remove it.

## Watermark Operations

### Set Watermark

Upload and set a watermark image.

**Reference:** [references/watermark-set.md](references/watermark-set.md)

**Common Tasks:**

- Set watermark: `yutu watermark set --file image.png --inVideoPosition bottomRight`

### Unset Watermark

Remove the watermark.

**Reference:** [references/watermark-unset.md](references/watermark-unset.md)

**Common Tasks:**

- Remove watermark: `yutu watermark unset`

## Resources

- [references/watermark-set.md](references/watermark-set.md): Set watermark.
- [references/watermark-unset.md](references/watermark-unset.md): Unset watermark.
