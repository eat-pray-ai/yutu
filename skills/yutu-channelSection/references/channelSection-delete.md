# ChannelSection Delete Command

Delete channel sections by ids.

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
# Delete a channel section
yutu channelSection delete --ids SECTION_ID
```
