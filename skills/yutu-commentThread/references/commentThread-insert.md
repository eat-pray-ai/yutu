# CommentThread Insert Command

Insert a new comment thread.

## Usage

```bash
yutu commentThread insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--authorChannelId` | `-a` | Channel id of the comment author |
| `--channelId` | `-c` | Channel id of the video owner |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--textOriginal` | `-t` | Text of the comment |
| `--videoId` | `-v` | ID of the video |

## Examples

```bash
# Insert a new comment thread on a video
yutu commentThread insert --videoId VIDEO_ID --textOriginal "Nice video!"
```
