---
name: youtube
description: "Use whenever the user mentions YouTube, video uploads, channel management, playlists, video SEO, or any YouTube Data API operation. Manages videos, playlists, comments, captions, subscriptions, thumbnails, analytics, and more."
license: MIT
metadata:
  author: eat-pray-ai
  homepage: "https://github.com/eat-pray-ai/yutu"
---

# YouTube

Manage YouTube resources via MCP tool calls — abuseReport, activity, caption, channel, channelBanner, channelSection, comment, commentThread, i18nLanguage, i18nRegion, liveBroadcast, liveChatBan, liveChatMessage, liveChatModerator, liveStream, member, membershipsLevel, playlist, playlistImage, playlistItem, search, subscription, superChatEvent, thirdPartyLink, thumbnail, video, videoAbuseReportReason, videoCategory, watermark.

## Key Principles

- **Confirm before deleting**: Always list exactly what will be deleted and get explicit user confirmation before calling any destructive tool.
- **Verify targets**: Before modifying or deleting, use a read tool to confirm the resource exists and matches the user's intent.
- **Preserve existing data**: When updating resources, only modify fields the user explicitly asked to change.
- **Report resource IDs**: Always include resource IDs in responses so the user can reference them.

## Tools

### Read (safe, no side effects)

- `activity-list`
- `caption-download`
- `caption-list`
- `channel-list`
- `channelSection-list`
- `comment-list`
- `commentThread-list`
- `i18nLanguage-list`
- `i18nRegion-list`
- `liveBroadcast-list`
- `liveChatMessage-list`
- `liveChatModerator-list`
- `liveStream-list`
- `member-list`
- `membershipsLevel-list`
- `playlist-list`
- `playlistImage-list`
- `playlistItem-list`
- `search-list`
- `subscription-list`
- `superChatEvent-list`
- `thirdPartyLink-list`
- `video-getRating`
- `video-list`
- `videoAbuseReportReason-list`
- `videoCategory-list`

### Write (create or update content)

- `abuseReport-insert`
- `caption-insert`
- `caption-update`
- `channel-update`
- `channelBanner-insert`
- `comment-insert`
- `comment-markAsSpam`
- `comment-setModerationStatus`
- `comment-update`
- `commentThread-insert`
- `liveBroadcast-bind`
- `liveBroadcast-insert`
- `liveBroadcast-insertCuepoint`
- `liveBroadcast-transition`
- `liveBroadcast-update`
- `liveChatBan-insert`
- `liveChatMessage-insert`
- `liveChatMessage-transition`
- `liveChatModerator-insert`
- `liveStream-insert`
- `liveStream-update`
- `playlist-insert`
- `playlist-update`
- `playlistImage-insert`
- `playlistImage-update`
- `playlistItem-insert`
- `playlistItem-update`
- `subscription-insert`
- `thirdPartyLink-insert`
- `thirdPartyLink-update`
- `thumbnail-set`
- `video-insert`
- `video-rate`
- `video-reportAbuse`
- `video-update`
- `watermark-set`

### Destructive (irreversible deletions — confirm with user first)

- `caption-delete`
- `channelSection-delete`
- `comment-delete`
- `liveBroadcast-delete`
- `liveChatBan-delete`
- `liveChatMessage-delete`
- `liveChatModerator-delete`
- `liveStream-delete`
- `playlist-delete`
- `playlistImage-delete`
- `playlistItem-delete`
- `subscription-delete`
- `thirdPartyLink-delete`
- `video-delete`
- `watermark-unset`

## Common Workflows

See [references/workflows.md](references/workflows.md) for step-by-step walkthroughs.

## YouTube Growth Tips

See [references/seo-guide.md](references/seo-guide.md) for the full guide. When creating or updating content, apply these principles:

- **Titles**: Curiosity gaps + power words. Front-load keywords. Under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Consistent schedule matters.
- **Engagement**: Pin a comment with a question. Reply within the first hour.
