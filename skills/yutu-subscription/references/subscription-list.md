# Subscription List Command

List subscriptions' info.

## Usage

```bash
yutu subscription list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--channelId` | `-c` | Return the subscriptions of the given channel owner |
| `--forChannelId` | `-C` | Return the subscriptions to the subset of these channels that the authenticated user is subscribed to |
| `--ids` | `-i` | Return the subscriptions with the given ids |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` | Return the subscriptions of the authenticated user (default true) |
| `--myRecentSubscribers` | `-R` | true\|false\|null |
| `--mySubscribers` | `-S` | Return the subscribers of the given channel owner |
| `--onBehalfOfContentOwner` | `-b` | |
| `--onBehalfOfContentOwnerChannel` | `-B` | |
| `--order` | `-O` | subscriptionOrderUnspecified\|relevance\|unread\|alphabetical (default "relevance") |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List my subscriptions
yutu subscription list --mine

# List my subscribers
yutu subscription list --mySubscribers
```
