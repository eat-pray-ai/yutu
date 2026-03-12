# Memberships Level List

List memberships levels. Use this skill to list information about channel membership levels.

## Usage

```bash
yutu membershipsLevel list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List membership levels in JSON format
yutu membershipsLevel list --output json
```
