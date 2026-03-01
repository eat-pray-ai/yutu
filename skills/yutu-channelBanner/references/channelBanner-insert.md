# ChannelBanner Insert Command

Insert a YouTube channel banner. Use this tool when you need to insert or upload a channel banner.

## Usage

```bash
yutu channelBanner insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | ID of the channel to insert the banner for |
| `--file` | `-f` | Path to the banner image |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` | YouTube channel ID linked to the content owner |
| `--output` | `-o` | json\|yaml\|silent |

## Examples

```bash
yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.jpg
yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.png --output json
```
