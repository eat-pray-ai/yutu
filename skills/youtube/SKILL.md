---
name: youtube
description: "Use whenever the user mentions YouTube, video uploads, channel management, playlists, video SEO, or any YouTube Data API operation. Manages videos, playlists, comments, captions, subscriptions, thumbnails, analytics, and more via the yutu CLI."
license: MIT
compatibility: Requires the yutu CLI binary (installable via npm, brew, or winget) and Google Cloud OAuth credentials for YouTube Data API v3.
metadata:
  author: eat-pray-ai
  homepage: "https://github.com/eat-pray-ai/yutu"
---

# YouTube

Manage YouTube resources using the `yutu` CLI — videos, playlists, comments, channels, captions, subscriptions, and more.

## Quick Start

1. Ensure `yutu` is installed and authenticated. If not, follow [references/setup.md](references/setup.md).
2. Identify the resource and operation from the tables below.
3. Run `yutu <resource> <operation> -h` for full flag details on any command.
4. For multistep tasks (upload + thumbnail + playlist), see [references/workflows.md](references/workflows.md).

## Key Principles

- Always verify before destructive operations — deletions are irreversible.
- Use `--output json` when you need to parse or chain results.
- Get your channel ID with `yutu channel list --for mine` — many operations need it.
- When updating metadata, only specify the fields you want to change.

## Operations

### Content

| Resource | Operations |
|----------|------------|
| caption | delete, download, insert, list, update |
| liveBroadcast | bind, delete, insert, insertCuepoint, list, transition, update |
| liveStream | delete, insert, list, update |
| thumbnail | set |
| video | delete, getRating, insert, list, rate, reportAbuse, update |
| watermark | set, unset |

### Organization

| Resource | Operations |
|----------|------------|
| playlist | delete, insert, list, update |
| playlistImage | delete, insert, list, update |
| playlistItem | delete, insert, list, update |

### Community

| Resource | Operations |
|----------|------------|
| abuseReport | insert |
| comment | delete, insert, list, markAsSpam, setModerationStatus, update |
| commentThread | insert, list |
| liveChatBan | delete, insert |
| liveChatMessage | delete, insert, list, transition |
| liveChatModerator | delete, insert, list |
| member | list |
| membershipsLevel | list |
| subscription | delete, insert, list |
| superChatEvent | list |

### Channel

| Resource | Operations |
|----------|------------|
| channel | list, update |
| channelBanner | insert |
| channelSection | delete, list |
| thirdPartyLink | delete, insert, list, update |

### Discovery

| Resource | Operations |
|----------|------------|
| activity | list |
| search | list |

### Metadata

| Resource | Operations |
|----------|------------|
| i18nLanguage | list |
| i18nRegion | list |
| videoAbuseReportReason | list |
| videoCategory | list |

## Common Workflows

See [references/workflows.md](references/workflows.md) for step-by-step walkthroughs of each task below.

| Task | Quick Command |
|------|---------------|
| Publishing video pipeline | `yutu video insert --privacy unlisted` → review → `yutu video update --privacy public` |
| Find unlisted/private videos | `yutu playlistItem list` (uploads playlist) → `yutu video list --parts id,snippet,status` |

## YouTube Growth Tips

See [references/seo-guide.md](references/seo-guide.md) for the full guide. When uploading or updating video metadata, apply these principles:

- **Titles**: Curiosity gaps + power words. Front-load keywords. Under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Consistent schedule matters.
- **Engagement**: Pin a comment with a question. Reply within the first hour.
