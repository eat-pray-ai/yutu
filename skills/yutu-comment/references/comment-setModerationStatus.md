# Comment SetModerationStatus Command

Set YouTube comments moderation status by ids.

## Usage

```bash
yutu comment setModerationStatus [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--banAuthor` | `-A` | If set to true the author of the comment gets added to the ban list |
| `--ids` | `-i` | IDs of comments |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--moderationStatus` | `-s` | heldForReview\|published\|rejected |
| `--output` | `-o` | json\|yaml\|silent |

## Examples

```bash
# Publish a comment held for review
yutu comment setModerationStatus --ids COMMENT_ID --moderationStatus published
```
