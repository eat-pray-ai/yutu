# Subscription List

List subscription information. Use this skill to list subscription information.

## Usage

```bash
yutu subscription list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--channelId` | `-c` |  | Return the subscriptions of the given channel owner |
| `--forChannelId` | `-C` |  | Return the subscriptions to the subset of these channels that the authenticated user is subscribed to |
| `--ids` | `-i` |  | Return the subscriptions with the given ids for Stubby or Apiary |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--maxResults` | `-n` |  | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` |  | Return the subscriptions of the authenticated user (default true) |
| `--myRecentSubscribers` | `-R` |  | true\|false\|null |
| `--mySubscribers` | `-S` |  | Return the subscribers of the given channel owner |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` |  | YouTube channel ID linked to the content owner |
| `--order` | `-O` |  | subscriptionOrderUnspecified\|relevance\|unread\|alphabetical (default "relevance") |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List my subscriptions
yutu subscription list --mine
# List subscriptions by channel ID with limit
yutu subscription list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10
# List specific subscriptions in JSON format
yutu subscription list --ids abc123,def456 --output json
# List subscriptions for a channel alphabetically
yutu subscription list --forChannelId UC_x5XG1OV2P6uZZ5FSM9Ttw --order alphabetical
```
