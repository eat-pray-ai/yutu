---
name: youtube
description: "Use when working with YouTube — upload videos, search content, manage playlists and channels, post and moderate comments, handle subscriptions and memberships, add captions, set thumbnails, check analytics, or any YouTube Data API operation via the yutu CLI."
metadata:
  openclaw:
    requires:
      env:
        - YUTU_CREDENTIAL
        - YUTU_CACHE_TOKEN
      bins:
        - yutu
      config:
        - client_secret.json
        - youtube.token.json
    primaryEnv: YUTU_CREDENTIAL
    emoji: "\U0001F3AC\U0001F430"
    homepage: https://github.com/eat-pray-ai/yutu
    install:
      - kind: node
        package: "@eat-pray-ai/yutu"
        bins: [yutu]
---

# YouTube

Manage YouTube resources using the yutu CLI — videos, playlists, comments, channels, captions, subscriptions, and more.

## Before You Begin

yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.

## Operations

Run `yutu <resource> <verb> -h` for full flag details and examples.

### Content

| Resource | Operation | Description |
|----------|-----------|-------------|
| caption | delete | Delete captions |
| caption | download | Download a caption |
| caption | insert | Insert a caption |
| caption | list | List captions |
| caption | update | Update a video caption |
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
| comment | delete | Delete comments |
| comment | insert | Create a comment |
| comment | list | List comments |
| comment | markAsSpam | Mark comments as spam |
| comment | setModerationStatus | Set comment moderation status |
| comment | update | Update a comment on a video |
| comment thread | insert | Insert a new comment thread |
| comment thread | list | List comment threads |
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

See [references/workflows.md](references/workflows.md) for detailed walkthroughs.

- **Upload a video**: `yutu video insert --file video.mp4 --title "..." --privacy public`, then optionally set thumbnail
- **Update video metadata**: Fetch current with `yutu video list --id VIDEO_ID`, then update changed fields
- **Create playlist + add videos**: Create with `yutu playlist insert`, find videos with `yutu search list --forMine`, add with `yutu playlistItem insert`
- **Post a comment**: Get channel ID with `yutu channel list --mine`, find video, then `yutu commentThread insert`
- **Channel analytics**: `yutu channel list --mine` + `yutu search list --forMine` + `yutu video list --id ...`
- **Competitor analysis**: `yutu channel list --forHandle @handle` + compare stats and top videos
- **Delete content**: Always verify with a list command first, then delete — deletions are irreversible
- **Subscribe/unsubscribe**: Check with `yutu subscription list --mine --forChannelId ...` before acting

## YouTube Growth Tips

See [references/seo-guide.md](references/seo-guide.md) for the full guide.

- **Titles**: Use curiosity gaps and power words. Front-load keywords. Keep under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, and 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Maintain consistent schedule.
- **Engagement**: Pin a comment with a question. Reply within the first hour.
