# Upload video to YouTube

A github action to upload video to YouTube using [yutu](https://eat-pray-ai/yutu).

## Example

```yaml
name: Upload video to YouTube
on:
  workflow_dispatch:

jobs:
  upload:
    runs-on: ubuntu-latest
    steps:
      - shell: bash
        run: |
          echo "generate video.mp4"
          ls ./*
          #
      - uses: eat-pray-ai/yutu/actions/upload@main
        with:
          token: ${{ secrets.YOUTUBE_TOKEN }}
          flags: "-f video.mp4 -t title -d description -g category -p public"
```
## Inputs

### `version`

**Optional**, The version of yutu to use, `v0.8.2` for example, leave empty to use the latest version.

### `token`

**Required**, Please refer to [yutu prerequisites](https://github.com/eat-pray-ai/yutu?tab=readme-ov-file#prerequisites) to generate a token.

### `flags`

**Required**, Flags passed to yutu, for all available flags when uploading a video:

```shell
❯ yutu video insert --help
upload a video to YouTube, with the specified title, description, tags, etc.

Usage:
  yutu video insert [flags]

Flags:
  -a, --autoLevels string                      true or false
  -g, --categoryId string                      Category of the video
  -c, --channelId string                       Channel ID of the video
  -d, --description string                     Description of the video
  -e, --embeddable                             Whether the video is embeddable (default true)
  -f, --file string                            Path to the video file
  -k, --forKids                                Whether the video is for kids
  -h, --help                                   help for insert
  -l, --language string                        Language of the video
  -L, --license string                         youtube(default) or creativeCommon (default "youtube")
  -n, --notifySubscribers                      Notify the channel subscribers about the new video (default true)
  -b, --onBehalfOfContentOwner string
  -B, --onBehalfOfContentOwnerChannel string
  -y, --playlistId string                      Playlist ID of the video
  -p, --privacy string                         Privacy status of the video: public, private, or unlisted
  -P, --publicStatsViewable                    Whether the extended video statistics can be viewed by everyone
  -U, --publishAt string                       Datetime when the video is scheduled to publish
  -s, --stabilize string                       true or false
  -T, --tags stringArray                       Comma separated tags
  -u, --thumbnail string                       Path to the thumbnail
  -t, --title string                           Title of the video
```
