# Comment Update Command

Update a comment on a video.

## Usage

```bash
yutu comment update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--canRate` | `-R` | Whether the viewer can rate the comment |
| `--id` | `-i` | ID of the comment |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--textOriginal` | `-t` | Text of the comment |
| `--viewerRating` | `-r` | none\|like\|dislike |

## Examples

```bash
# Update comment text
yutu comment update --id COMMENT_ID --textOriginal "Updated text"
```
