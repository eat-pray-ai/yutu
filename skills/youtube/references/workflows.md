# YouTube Workflows

End-to-end workflows for multi-step YouTube tasks using the `yutu` CLI.

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

# 4. Optionally change privacy to public
yutu video update --id VIDEO_ID --privacy public
```

Key points:
- `search list` and `search list --channelId` only return **public** videos.
- The uploads playlist ID follows the pattern: channel `UCxxx` → playlist `UUxxx`.
- `video list --ids` accepts at most 50 IDs per call — batch accordingly.
- The `status` part contains `privacyStatus` (public/unlisted/private).
