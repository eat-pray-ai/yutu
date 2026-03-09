# Comment Insert

Create a comment. Use this skill to create a comment on a video.

## Usage

```bash
yutu comment insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--authorChannelId` | `-a` |  | Channel id of the comment author |
| `--canRate` | `-R` |  | Whether the viewer can rate the comment |
| `--channelId` | `-c` |  | Channel id of the video owner |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--parentId` | `-P` | Yes | ID of the parent comment |
| `--textOriginal` | `-t` | Yes | Text of the comment |
| `--videoId` | `-v` |  | ID of the video |

## Examples

```bash
# Reply to a comment
yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Hello'
# Reply with rating enabled
yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Reply' --canRate
```
