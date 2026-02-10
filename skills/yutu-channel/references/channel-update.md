# Channel Update Command

Update channel's info, such as title, description, etc

## Usage

```bash
yutu channel update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--country` | `-c` | Country of the channel |
| `--customUrl` | `-u` | Custom URL of the channel |
| `--defaultLanguage` | `-l` | The language of the channel's default title and description |
| `--description` | `-d` | Description of the channel |
| `--id` | `-i` | ID of the channel to update |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |
| `--title` | `-t` | Title of the channel |

## Examples

```bash
# Update channel title
yutu channel update --id UCxxxxxx --title "New Title"

# Update channel description
yutu channel update --id UCxxxxxx --description "New Description"
```
