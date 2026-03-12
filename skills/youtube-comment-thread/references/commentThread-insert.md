# Comment Thread Insert

Insert a new comment thread. Use this skill to insert a new comment thread.

## Usage

```bash
yutu commentThread insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--authorChannelId` | `-a` |  | Channel id of the comment author |
| `--channelId` | `-c` | Yes | Channel id of the video owner |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--textOriginal` | `-t` | Yes | Text of the comment |
| `--videoId` | `-v` | Yes | ID of the video |

## Examples

```bash
# Post a comment on a video
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Great video!'
# Post a comment with JSON output
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Nice work!' --output json
```
