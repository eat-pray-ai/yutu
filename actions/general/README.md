# yutu

A github action to use [yutu](https://eat-pray-ai/yutu).

## Example

See [actions/upload](../upload/README.md) for an example.

## Inputs

### `version`

**Optional**, The version of yutu to use, `v0.8.2` for example, leave empty to use the latest version.

### `token`

**Required**, Please refer to [yutu prerequisites](https://github.com/eat-pray-ai/yutu?tab=readme-ov-file#prerequisites) to generate a token.

### `command`

**Optional**, The command to run, for all available commands:

```shell
❯ yutu --help
yutu is a fully functional CLI for YouTube, which can be used to manupulate YouTube videos, playlists, channels, etc.

Usage:
  yutu [flags]
  yutu [command]

Available Commands:
  activity               list YouTube activities
  auth                   authenticate with YouTube API
  channel                manipulate YouTube channels
  channelBanner          insert Youtube channelBanner
  completion             Generate the autocompletion script for the specified shell
  help                   Help about any command
  i18nLanguage           list YouTube i18nLanguages
  i18nRegion             list YouTube i18nRegions
  member                 list YouTube members
  membershipsLevel       list YouTube memberships levels
  playlist               manipulate YouTube playlists
  playlistItem           manipulate YouTube playlist items
  search                 Search for youtube resources
  subscription           manipulate YouTube subscriptions
  version                Show the version of yutu
  video                  manipulate YouTube videos
  videoAbuseReportReason list YouTube video abuse report reasons
  videoCategory          list YouTube video categories
  watermark              manipulate Youtube watermarks

Flags:
  -h, --help   help for yutu

Use "yutu [command] --help" for more information about a command.
```

### `subcommand`

**Optional**, The subcommand corresponding to the command, to get all available subcommands for a command:

```shell
❯ yutu <command> --help
❯ yutu playlist --help                                                                            
manipulate YouTube playlists, such as insert, update, etc.

Usage:
  yutu playlist [flags]
  yutu playlist [command]

Available Commands:
  delete      delete a playlist
  insert      create a new playlist
  list        list playlist's info
  update      update an existing playlist

Flags:
  -h, --help   help for playlist

Use "yutu playlist [command] --help" for more information about a command.
```

### `flags`

**Optional**, Flags passed to yutu, to get all available flags for a command and subcommand:

```shell
❯ yutu <command> <subcommand> -h
❯ yutu subscription insert --help                                                 
Insert a subscription

Usage:
  yutu subscription insert [flags]

Flags:
  -c, --channelId string             ID of the channel to be subscribed
  -d, --description string           Description of the subscription
  -h, --help                         help for insert
  -s, --subscriberChannelId string   Subscriber's channel ID
  -t, --title string                 Title of the subscription
```

