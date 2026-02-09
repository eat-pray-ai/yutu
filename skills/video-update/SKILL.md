---
name: video-update
description: Update video metadata.
---

# Video Update

This skill provides instructions for updating YouTube videos using the `yutu` CLI.

## Usage

```bash
yutu video update [flags]
```

## Options

- `--id`, `-i`: ID of the video to update.
- `--title`, `-t`: Title of the video.
- `--description`, `-d`: Description of the video.
- `--tags`, `-a`: Comma separated tags.
- `--categoryId`, `-g`: Category of the video.
- `--privacy`, `-p`: Privacy status: `public`, `private`, `unlisted`.
- `--thumbnail`, `-u`: Path to the thumbnail file.
- `--language`, `-l`: Language of the video.
- `--license`, `-L`: `youtube` or `creativeCommon`.
- `--embeddable`, `-E`: Whether the video is embeddable.
- `--playlistId`, `-y`: Update video's playlist.
- `--output`, `-o`: Output format (`json`, `yaml`, `silent`).
- `--jsonPath`, `-j`: JSONPath expression to filter the output.

## Examples

**Update video title:**

```bash
yutu video update --id VIDEO_ID --title "New Title"
```
