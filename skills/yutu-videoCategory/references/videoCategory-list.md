# VideoCategory List Command

List YouTube video categories' info, such as id, title, assignable, etc.

## Usage

```bash
yutu videoCategory list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--hl` | `-l` | Host language |
| `--ids` | `-i` | IDs of the video categories |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |
| `--regionCode` | `-r` | Region code (default "US") |

## Examples

```bash
# List video categories for US
yutu videoCategory list --regionCode US

# List categories for a specific language
yutu videoCategory list --hl es
```
