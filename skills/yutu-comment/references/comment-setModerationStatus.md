# Comment SetModerationStatus Command

Set comment moderation status. Use this tool when you need to set comment moderation status.

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
# Publish a held comment
yutu comment setModerationStatus --ids abc123 --moderationStatus published
# Hold multiple comments for review
yutu comment setModerationStatus --ids abc123,def456 --moderationStatus heldForReview
# Reject a comment and ban author
yutu comment setModerationStatus --ids abc123 --moderationStatus rejected --banAuthor
```
