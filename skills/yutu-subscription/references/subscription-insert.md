# Subscription Insert Command

Insert a YouTube subscription. Use this tool when you need to insert a YouTube subscription.

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
yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5X --title 'Google Developers'
```
