# CommentThread Insert Command

Insert a new comment thread. Use this tool when you need to insert a new comment thread.

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
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Great video!'
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Nice work!' --output json
```
