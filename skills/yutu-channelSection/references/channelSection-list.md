# Channel Section List

List channel sections. Use this skill to list channel sections.

## Usage

```bash
yutu channelSection list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` |  | Return the ChannelSections owned by the specified channel id |
| `--hl` | `-l` |  | Return content in specified language |
| `--ids` | `-i` |  | Return the channel sections with the given ids |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--mine` | `-M` |  | Return the ChannelSections owned by the authenticated user |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List my channel sections
yutu channelSection list --mine
# List channel sections by channel ID
yutu channelSection list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
# List specific channel sections in JSON format
yutu channelSection list --ids abc123,def456 --output json
```
