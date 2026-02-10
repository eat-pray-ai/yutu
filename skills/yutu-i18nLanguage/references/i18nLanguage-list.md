# I18nLanguage List Command

List YouTube i18n languages' id, hl, and name.

## Usage

```bash
yutu i18nLanguage list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--hl` | `-l` | Host language |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List supported languages
yutu i18nLanguage list

# List languages in a specific host language
yutu i18nLanguage list --hl es
```
