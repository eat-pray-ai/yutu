# Channel Update

Update channel information. Use this skill to update channel information.

## Usage

```bash
yutu channel update [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--country` | `-c` |  | Country of the channel |
| `--customUrl` | `-u` |  | Custom URL of the channel |
| `--defaultLanguage` | `-l` |  | The language of the channel's default title and description |
| `--description` | `-d` |  | Description of the channel |
| `--id` | `-i` | Yes | ID of the channel to update |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--title` | `-t` |  | Title of the channel |

## Examples

```bash
# Update channel description
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --description 'New description'
# Update channel title and country
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --title 'New Title' --country US
# Update channel default language
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --defaultLanguage en
```
