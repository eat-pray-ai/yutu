# Comment Insert Command

Insert a comment to a video.

## Usage

```bash
yutu comment insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--authorChannelId` | `-a` | Channel id of the comment author |
| `--canRate` | `-R` | Whether the viewer can rate the comment |
| `--channelId` | `-c` | Channel id of the video owner |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--parentId` | `-P` | ID of the parent comment |
| `--textOriginal` | `-t` | Text of the comment |
| `--videoId` | `-v` | ID of the video |

## Examples

```bash
# Comment on a video
yutu comment insert --videoId VIDEO_ID --textOriginal "Great video!"

# Reply to a comment
yutu comment insert --parentId PARENT_COMMENT_ID --textOriginal "I agree!"
```
