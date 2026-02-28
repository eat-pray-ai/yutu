# ChannelBanner Insert Command

Insert Youtube channel banner.

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
# Insert banner for my channel
yutu channelBanner insert --file banner.jpg

# Insert banner for a specific channel
yutu channelBanner insert --channelId CHANNEL_ID --file banner.jpg
```
