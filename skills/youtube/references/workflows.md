# YouTube Workflows

Step-by-step workflows for common YouTube tasks using the `yutu` CLI.

## Content Management

### Upload a video

```bash
yutu video insert --file video.mp4 --title 'My Video' --categoryId 22 --privacy public

# With optional fields
yutu video insert --file video.mp4 --title 'Tutorial' --categoryId 27 --privacy public \
  --description 'Learn Go in 10 minutes' --tags 'go,tutorial' --thumbnail cover.jpg
```

### Update video metadata (title, description, tags)

```bash
# Get current metadata to see existing values
yutu video list --ids VIDEO_ID --output json

# Update only the fields you want to change (preserve all others)
yutu video update --id VIDEO_ID --title 'New Title'
yutu video update --id VIDEO_ID --description 'Updated description' --privacy public
yutu video update --id VIDEO_ID --tags 'music,pop,2024' --categoryId 10
```

### Set a video thumbnail

```bash
# Verify the video exists
yutu video list --ids VIDEO_ID

# Set the thumbnail
yutu thumbnail set --file image.jpg --videoId VIDEO_ID
```

### Update captions/subtitles

```bash
# List existing captions to find IDs and languages
yutu caption list --videoId VIDEO_ID

# Update an existing caption or insert a new one
yutu caption update --videoId VIDEO_ID --language en --name English
yutu caption insert --file subtitle.srt --videoId VIDEO_ID --language en --name English
```

## Channel Management

### Update channel branding

```bash
# Get your channel info
yutu channel list --mine

# Update channel metadata
yutu channel update --id CHANNEL_ID --description 'New description'
yutu channel update --id CHANNEL_ID --title 'New Title' --country US
yutu channel update --id CHANNEL_ID --defaultLanguage en

# Update the channel banner image
yutu channelBanner insert --channelId CHANNEL_ID --file banner.jpg
```

## Playlists

### Create a playlist and add videos

```bash
# Get your channel ID
yutu channel list --mine

# Create the playlist
yutu playlist insert --title 'My Playlist' --channelId CHANNEL_ID --privacy public

# Find videos to add (e.g., your 5 most recent)
yutu search list --forMine --order date --types video --maxResults 5

# Add each video to the playlist
yutu playlistItem insert --kind video --playlistId PLAYLIST_ID --channelId CHANNEL_ID --kVideoId VIDEO_ID
```

### Add videos to an existing playlist

```bash
# Check what is already in the playlist
yutu playlistItem list --playlistId PLAYLIST_ID

# Add new videos
yutu playlistItem insert --kind video --playlistId PLAYLIST_ID --channelId CHANNEL_ID --kVideoId VIDEO_ID
```

## Comments & Community

### Post a comment on a video

```bash
# Get your channel ID
yutu channel list --mine

# Find the video (e.g., your latest)
yutu search list --forMine --order date --types video --maxResults 1

# Post the comment
yutu commentThread insert --channelId CHANNEL_ID --videoId VIDEO_ID \
  --authorChannelId CHANNEL_ID --textOriginal 'Great video!'
```

### Reply to a comment

```bash
# Find the comment thread
yutu commentThread list --videoId VIDEO_ID --maxResults 10

# Post the reply
yutu comment insert --channelId CHANNEL_ID --videoId VIDEO_ID \
  --authorChannelId CHANNEL_ID --parentId THREAD_ID --textOriginal 'Thanks!'
```

### Subscribe to a channel

```bash
# Look up the target channel
yutu channel list --forHandle @TargetChannel

# Check if already subscribed
yutu subscription list --mine --forChannelId TARGET_CHANNEL_ID

# Subscribe if not already subscribed
yutu subscription insert --subscriberChannelId MY_CHANNEL_ID --channelId TARGET_CHANNEL_ID
```

## Destructive Operations

> These operations are irreversible. Always verify targets before executing.

### Delete a video

```bash
# Verify the video exists and confirm its title
yutu video list --ids VIDEO_ID

# Delete the video
yutu video delete --ids VIDEO_ID
```

### Delete a video by criteria (e.g., oldest video)

```bash
# Search for the target video
yutu search list --forMine --order date --types video --maxResults 1

# Verify and delete
yutu video list --ids VIDEO_ID
yutu video delete --ids VIDEO_ID
```

### Delete comments

```bash
# Find matching comments
yutu commentThread list --videoId VIDEO_ID --searchTerms 'spam'

# Delete each matching comment
yutu comment delete --ids COMMENT_ID
```

### Delete a playlist

```bash
# Verify the playlist
yutu playlist list --ids PLAYLIST_ID

# Delete
yutu playlist delete --ids PLAYLIST_ID
```

### Remove items from a playlist

```bash
# List items and get their playlistItem IDs (not video IDs)
yutu playlistItem list --playlistId PLAYLIST_ID

# Delete items
yutu playlistItem delete --ids PLAYLIST_ITEM_ID
```

### Unsubscribe from a channel

```bash
# Find the subscription ID
yutu subscription list --mine --forChannelId TARGET_CHANNEL_ID

# Delete the subscription
yutu subscription delete --ids SUBSCRIPTION_ID
```

### Delete captions

```bash
# List captions to find IDs
yutu caption list --videoId VIDEO_ID

# Delete
yutu caption delete --ids CAPTION_ID
```

### Remove a watermark

```bash
# Get the channel ID
yutu channel list --mine

# Unset the watermark
yutu watermark unset --channelId CHANNEL_ID
```

## Analytics & Research

### Channel performance

```bash
# Get channel stats
yutu channel list --mine --output json

# List recent videos
yutu search list --forMine --order date --types video --maxResults 10

# Get detailed stats for those videos
yutu video list --ids VIDEO_ID_1,VIDEO_ID_2,VIDEO_ID_3 --output json

# Analyze: subscriber count, total views, per-video performance,
# engagement rates, upload frequency.
```

### Video performance

```bash
# Get detailed video info
yutu video list --ids VIDEO_ID --output json

# Sample engagement via comments
yutu commentThread list --videoId VIDEO_ID --maxResults 20

# Analyze: views, likes, comments, publish date, engagement rate, growth trends.
```

### Competitor analysis

```bash
# Look up the competitor channel
yutu channel list --forHandle @Competitor --output json

# Get your own channel for comparison
yutu channel list --mine --output json

# Find the competitor's top videos
yutu search list --channelId COMPETITOR_CHANNEL_ID --order viewCount --types video --maxResults 10

# Compare: subscribers, upload frequency, view counts, content themes.
```

### Recent activity

```bash
# Get your channel info
yutu channel list --mine

# List recent activity (uploads, likes, favorites, comments)
yutu activity list --channelId CHANNEL_ID --maxResults 20
```
