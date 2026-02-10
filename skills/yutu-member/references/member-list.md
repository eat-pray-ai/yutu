# Member List Command

List channel's members' info, such as channelId, displayName, etc.

## Usage

```bash
yutu member list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--hasAccessToLevel` | `-a` | Filter members in the results set to the ones that have access to a level |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--maxResults` | `-n` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--memberChannelId` | `-c` | Comma separated list of channel Ids. Only data about members that are part of this list will be included |
| `--mode` | `-m` | listMembersModeUnknown, updates, or all_current (default "all_current") |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [snippet]) |

## Examples

```bash
# List all current members
yutu member list

# List members with access to a specific level
yutu member list --hasAccessToLevel LEVEL_ID
```
