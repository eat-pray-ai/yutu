# I18nRegion List Command

List YouTube i18n regions' id, hl, and name.

## Usage

```bash
yutu i18nRegion list [flags]
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
# List supported regions
yutu i18nRegion list

# List regions for a specific host language
yutu i18nRegion list --hl es
```
