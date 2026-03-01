# MembershipsLevel List Command

List memberships levels. Use this tool when you need to list information about channel membership levels.

## Usage

```bash
yutu membershipsLevel list [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|table (default "table") |
| `--parts` | `-p` | Comma separated parts (default [id,snippet]) |

## Examples

```bash
yutu membershipsLevel list --output json
```
