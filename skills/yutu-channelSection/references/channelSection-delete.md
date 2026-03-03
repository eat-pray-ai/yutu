# ChannelSection Delete Command

Delete channel sections.

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
# Delete a channel section by ID
yutu channelSection delete --ids abc123
# Delete multiple channel sections
yutu channelSection delete --ids abc123,def456
```
