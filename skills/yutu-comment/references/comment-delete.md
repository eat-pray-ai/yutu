# Comment Delete Command

Delete comments. Use this tool when you need to delete comments by IDs.

## Usage

```bash
yutu comment delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of comments |

## Examples

```bash
# Delete a comment by ID
yutu comment delete --ids abc123
# Delete multiple comments
yutu comment delete --ids abc123,def456
```
