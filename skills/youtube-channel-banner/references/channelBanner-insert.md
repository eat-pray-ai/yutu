# Channel Banner Insert

Insert a channel banner. Use this skill to upload a channel banner.

## Usage

```bash
yutu channelBanner insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` | Yes | ID of the channel to insert the banner for |
| `--file` | `-f` | Yes | Path to the banner image |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` |  | YouTube channel ID linked to the content owner |
| `--output` | `-o` |  | json\|yaml\|silent |

## Examples

```bash
# Upload a channel banner
yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.jpg
# Upload a channel banner with JSON output
yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.png --output json
```
