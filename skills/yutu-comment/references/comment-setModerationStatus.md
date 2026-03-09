# Comment SetModerationStatus

Set comment moderation status. Use this skill to set comment moderation status.

## Usage

```bash
yutu comment setModerationStatus [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--banAuthor` | `-A` |  | If set to true the author of the comment gets added to the ban list |
| `--ids` | `-i` | Yes | IDs of comments |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--moderationStatus` | `-s` | Yes | heldForReview\|published\|rejected |
| `--output` | `-o` |  | json\|yaml\|silent |

## Examples

```bash
# Publish a held comment
yutu comment setModerationStatus --ids abc123 --moderationStatus published
# Hold multiple comments for review
yutu comment setModerationStatus --ids abc123,def456 --moderationStatus heldForReview
# Reject a comment and ban author
yutu comment setModerationStatus --ids abc123 --moderationStatus rejected --banAuthor
```
