# Subscription Insert Command

Insert a YouTube subscription.

## Usage

```bash
yutu subscription insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | ID of the channel to be subscribed |
| `--description` | `-d` | Description of the subscription |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--subscriberChannelId` | `-s` | Subscriber's channel id |
| `--title` | `-t` | Title of the subscription |

## Examples

```bash
# Subscribe to a channel
yutu subscription insert --channelId CHANNEL_ID
```
