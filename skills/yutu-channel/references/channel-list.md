# Channel List Command

List channel's info, such as title, description, etc.

## Usage

```bash
yutu channel list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--categoryId` | `-g` | Return the channels within the specified guide category id |
| `--forHandle` | `-d` | Return the channel associated with a YouTube handle |
| `--forUsername` | `-u` | Return the channel associated with a YouTube username |
| `--hl` | `-l` | Specifies the localization language of the metadata |
| `--ids` | `-i` | Return the channels with the specified Ids |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--managedByMe` | `-E` | Return the channels managed by the authenticated user |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--mine` | `-M` | Return the ids of channels owned by the authenticated user (default true) |
| `--mySubscribers` | `-S` | Return the channels subscribed to the authenticated user |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet,status]) |

## Examples

```bash
# List my channel
yutu channel list --mine

# List channel by ID
yutu channel list --ids UCxxxxxxxx

# List channel by handle
yutu channel list --forHandle @handle
```
