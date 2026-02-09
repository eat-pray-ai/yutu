---
name: video-insert
description: Upload a video to YouTube.
---

# Video Insert

This skill provides instructions for uploading YouTube videos using the `yutu` CLI.

## Usage

```bash
yutu video insert [flags]
```

## Options

- `--file`, `-f`: Path to the video file.
- `--title`, `-t`: Title of the video.
- `--description`, `-d`: Description of the video.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--tags`, `-a`: Comma separated tags.
- `--categoryId`, `-g`: Category of the video.
- `--language`, `-l`: Language of the video.
- `--playlistId`, `-y`: Add video to this playlist.
- `--thumbnail`, `-u`: Path to the thumbnail file.
- `--forKids`, `-K`: Whether the video is for kids.
- `--notifySubscribers`, `-N`: Notify subscribers (default true).
- `--embeddable`, `-E`: Whether video is embeddable (default true).
- `--license`, `-L`: `youtube` or `creativeCommon`.
- `--publicStatsViewable`, `-P`: Public stats viewable.
- `--publishAt`, `-U`: Scheduled publish time.
- `--channelId`, `-c`: Upload to this channel ID.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Upload a video:**

```bash
yutu video insert --file video.mp4 --title "My Video" --privacy public
```
