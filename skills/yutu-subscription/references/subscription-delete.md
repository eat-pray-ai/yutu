# Subscription Delete Command

Delete subscriptions. Use this tool when you need to delete subscriptions by IDs.

## Usage

```bash
yutu subscription delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the subscriptions to delete |

## Examples

```bash
yutu subscription delete --ids abc123
yutu subscription delete --ids abc123,def456
```
