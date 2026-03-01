# ChannelSection Delete Command

Delete channel sections. Use this tool when you need to delete channel sections by IDs.

## Usage

```bash
yutu channelSection delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | Delete the channel sections with the given ids |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |

## Examples

```bash
yutu channelSection delete --ids abc123
yutu channelSection delete --ids abc123,def456
```
