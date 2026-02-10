# CommentThread List Command

List YouTube comment threads.

## Usage

```bash
yutu commentThread list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--allThreadsRelatedToChannelId` | `-a` | Returns the comment threads of all videos of the channel and the channel comments as well |
| `--channelId` | `-c` | Channel id of the video owner |
| `--ids` | `-i` | IDs of the comment threads |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--moderationStatus` | `-m` | published\|heldForReview\|likelySpam\|rejected (default "published") |
| `--order` | `-O` | orderUnspecified\|time\|relevance (default "time") |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |
| `--searchTerms` | `-s` | Search terms |
| `--textFormat` | `-t` | textFormatUnspecified\|html (default "html") |
| `--videoId` | `-v` | Returns the comment threads of the specified video |

## Examples

```bash
# List comment threads for a video
yutu commentThread list --videoId VIDEO_ID
```
