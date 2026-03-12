# Comment Delete

Delete comments. Use this skill to delete comments by IDs.

## Usage

```bash
yutu comment delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of comments |

## Examples

```bash
# Delete a comment by ID
yutu comment delete --ids abc123
# Delete multiple comments
yutu comment delete --ids abc123,def456
```
