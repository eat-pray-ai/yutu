# Subscription Insert

Insert a new subscription. Use this skill to insert a new subscription.

## Usage

```bash
yutu subscription insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` | Yes | ID of the channel to be subscribed |
| `--description` | `-d` |  | Description of the subscription |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--subscriberChannelId` | `-s` | Yes | Subscriber's channel id |
| `--title` | `-t` |  | Title of the subscription |

## Examples

```bash
# Subscribe to a channel
yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
# Subscribe with a title
yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5X --title 'Google Developers'
```
