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
- Get your channel ID with `yutu channel list --mine` — many operations need it.
- When updating metadata, only specify the fields you want to change.

## Operations

### Content

| Resource | Operation | Description |
|----------|-----------|-------------|
| caption | delete | Delete captions |
| caption | download | Download a caption |
| caption | insert | Insert a caption |
| caption | list | List captions |
| caption | update | Update a video caption |
| live broadcast | bind | Bind a live broadcast to a stream |
| live broadcast | delete | Delete live broadcasts |
| live broadcast | insert | Insert a live broadcast |
| live broadcast | insertCuepoint | Insert a cuepoint into a live broadcast |
| live broadcast | list | List live broadcasts |
| live broadcast | transition | Transition a live broadcast |
| live broadcast | update | Update a live broadcast |
| live stream | delete | Delete live streams |
| live stream | insert | Insert a live stream |
| live stream | list | List live streams |
| live stream | update | Update a live stream |
| thumbnail | set | Set a thumbnail for a video |
| video | delete | Delete videos |
| video | getRating | Get video ratings |
| video | insert | Upload a video |
| video | list | List video information |
| video | rate | Rate a video |
| video | reportAbuse | Report abuse on a video |
| video | update | Update a video |
| watermark | set | Set a watermark for channel's videos |
| watermark | unset | Unset a watermark for channel's videos |

### Organization

| Resource | Operation | Description |
|----------|-----------|-------------|
| playlist | delete | Delete playlists |
| playlist | insert | Create a new playlist |
| playlist | list | List playlist information |
| playlist | update | Update a playlist |
| playlist image | delete | Delete playlist images |
| playlist image | insert | Insert a playlist image |
| playlist image | list | List playlist images |
| playlist image | update | Update a playlist image |
| playlist item | delete | Delete items from a playlist |
| playlist item | insert | Insert a playlist item into a playlist |
| playlist item | list | List playlist items |
| playlist item | update | Update a playlist item |

### Community

| Resource | Operation | Description |
|----------|-----------|-------------|
| abuse report | insert | Insert an abuse report |
| comment | delete | Delete comments |
| comment | insert | Create a comment |
| comment | list | List comments |
| comment | markAsSpam | Mark comments as spam |
| comment | setModerationStatus | Set comment moderation status |
| comment | update | Update a comment on a video |
| comment thread | insert | Insert a new comment thread |
| comment thread | list | List comment threads |
| live chat ban | delete | Delete live chat bans |
| live chat ban | insert | Insert a live chat ban |
| live chat message | delete | Delete live chat messages |
| live chat message | insert | Send a live chat message |
| live chat message | list | List live chat messages |
| live chat message | transition | Transition a live chat message |
| live chat moderator | delete | Delete live chat moderators |
| live chat moderator | insert | Insert a live chat moderator |
| live chat moderator | list | List live chat moderators |
| member | list | List channel members |
| memberships level | list | List memberships levels |
| subscription | delete | Delete subscriptions |
| subscription | insert | Insert a new subscription |
| subscription | list | List subscription information |
| super chat event | list | List Super Chat events |

### Channel

| Resource | Operation | Description |
|----------|-----------|-------------|
| channel | list | List channel information |
| channel | update | Update channel information |
| channel banner | insert | Insert a channel banner |
| channel section | delete | Delete channel sections |
| channel section | list | List channel sections |
| third party link | delete | Delete a third-party link |
| third party link | insert | Insert a new third-party link |
| third party link | list | List third-party links |
| third party link | update | Update a third-party link |

### Discovery

| Resource | Operation | Description |
|----------|-----------|-------------|
| activity | list | List activities |
| search | list | Search resources |

### Metadata

| Resource | Operation | Description |
|----------|-----------|-------------|
| i18n language | list | List i18n languages |
| i18n region | list | List i18n regions |
| video abuse report reason | list | List video abuse report reasons |
| video category | list | List video categories |

## Common Workflows

See [references/workflows.md](references/workflows.md) for step-by-step walkthroughs of each task below.

| Task | Quick Command |
|------|---------------|
| Upload a video | `yutu video insert --file video.mp4 --title "..." --privacy public` |
| Update video metadata | `yutu video list --ids ID` then `yutu video update --id ID --title "..."` |
| Create playlist + add videos | `yutu playlist insert` → `yutu playlistItem insert` |
| Post a comment | `yutu commentThread insert --channelId ... --videoId ... --textOriginal "..."` |
| Channel analytics | `yutu channel list --mine --output json` |
| Competitor analysis | `yutu channel list --forHandle @handle --output json` |
| Delete content | Always `list` first, then `delete` — irreversible |
| Subscribe/unsubscribe | Check `yutu subscription list --mine --forChannelId ...` before acting |

## YouTube Growth Tips

See [references/seo-guide.md](references/seo-guide.md) for the full guide. When uploading or updating video metadata, apply these principles:

- **Titles**: Curiosity gaps + power words. Front-load keywords. Under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Consistent schedule matters.
- **Engagement**: Pin a comment with a question. Reply within the first hour.
