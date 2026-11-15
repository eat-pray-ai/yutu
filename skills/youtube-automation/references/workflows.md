# YouTube Workflows

End-to-end workflows for managing YouTube content using the `yutu` CLI.

## Publishing Video Pipeline

Videos go through a staging queue before publishing. The "Unlisted" playlist acts as an inbox.

### 1. Upload (stage to queue)

Upload a video as unlisted. Use the `--description` field to describe the video content for AI context (what happens, tone, subjects) — this will be rewritten for SEO before publishing.

```bash
CHANNEL_ID=$(yutu channel list --for mine --output json | jq -r '.[0].id')

VIDEO_ID=$(yutu video insert --file video.mp4 \
  --title 'Draft title' \
  --description 'A husky finds an octopus in shallow water. Every time it gets close, the octopus sprays water at it. Funny and playful tone.' \
  --categoryId 22 --privacy unlisted --yes \
  --output json | jq -r '.id')

yutu playlistItem insert --kind video \
  --playlistId UNLISTED_PLAYLIST_ID --channelId "$CHANNEL_ID" \
  --kVideoId "$VIDEO_ID" --yes
```

### 2. Review (AI reads description, generates SEO fields)

Read the staged video's description to understand its content, then generate optimized title, description, tags, and thumbnail.

```bash
yutu video list --ids VIDEO_ID --parts id,snippet,status --output json
```

The AI uses the description as a content brief to generate:
- SEO title (curiosity gap + keywords, under 60 chars)
- SEO description (keywords, timestamps, CTAs, hashtags)
- Tags (broad + long-tail mix)
- Thumbnail guidance

### 3. Publish (update fields + make public + move to target playlist)

Overwrite all fields with SEO-optimized versions, change privacy to public, add to a themed playlist, and remove from the "Unlisted" queue.

```bash
CHANNEL_ID=$(yutu channel list --for mine --output json | jq -r '.[0].id')

# Update video metadata and publish
yutu video update --id VIDEO_ID \
  --title 'Husky vs Octopus Water Fight!' \
  --description 'Watch this hilarious husky get sprayed by an octopus...' \
  --tags 'husky,octopus,funny animals,pets' \
  --privacy public --yes

# Set thumbnail
yutu thumbnail set --file thumbnail.jpg --videoId VIDEO_ID

# Add to target playlist
yutu playlistItem insert --kind video \
  --playlistId TARGET_PLAYLIST_ID --channelId "$CHANNEL_ID" \
  --kVideoId VIDEO_ID --yes

# Remove from "Unlisted" queue (find the playlistItem ID first)
ITEM_ID=$(yutu playlistItem list --playlistId UNLISTED_PLAYLIST_ID --maxResults 0 --output json \
  | jq -r '.[] | select(.snippet.resourceId.videoId == "VIDEO_ID") | .id')
yutu playlistItem delete --ids "$ITEM_ID" --yes
```

## Find Unlisted/Private Videos

The search API only returns public videos. To find unlisted or private videos, use the uploads playlist which contains ALL uploads regardless of privacy status.

```bash
# 1. Get the uploads playlist ID (replace UC→UU in channel ID)
UPLOADS=$(yutu channel list --for mine --parts id,contentDetails --output json \
  | jq -r '.[0].contentDetails.relatedPlaylists.uploads')

# 2. List all videos in the uploads playlist
yutu playlistItem list --playlistId "$UPLOADS" --maxResults 0 --output json > uploads.json

# 3. Extract video IDs and query privacy status in batches of 50
cat uploads.json | jq -r '.[].snippet.resourceId.videoId' | \
  xargs -n50 | while read BATCH; do
    yutu video list --ids "$(echo $BATCH | tr ' ' ',')" \
      --parts id,snippet,status --output json
  done | jq '.[] | select(.status.privacyStatus == "unlisted") | {id, title: .snippet.title, published: .snippet.publishedAt}'
```

Key points:
- `search list` only returns **public** videos.
- The uploads playlist ID follows the pattern: channel `UCxxx` → playlist `UUxxx`.
- `video list --ids` accepts at most 50 IDs per call — batch accordingly.
- The `status` part contains `privacyStatus` (public/unlisted/private).
