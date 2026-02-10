# ChannelSection List Command

List channel sections.

## Usage

```bash
yutu channelSection list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | Return the ChannelSections owned by the specified channel id |
| `--hl` | `-l` | Return content in specified language |
| `--ids` | `-i` | Return the channel sections with the given ids |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--mine` | `-M` | Return the ChannelSections owned by the authenticated user |
| `--onBehalfOfContentOwner` | `-b` | |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List my channel sections
yutu channelSection list --mine

# List channel sections by ID
yutu channelSection list --ids SECTION_ID
```
