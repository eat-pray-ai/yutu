# Subscription Delete

Delete subscriptions. Use this skill to delete subscriptions by IDs.

## Usage

```bash
yutu subscription delete [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--ids` | `-i` | Yes | IDs of the subscriptions to delete |

## Examples

```bash
# Delete a subscription by ID
yutu subscription delete --ids abc123
# Delete multiple subscriptions
yutu subscription delete --ids abc123,def456
```
