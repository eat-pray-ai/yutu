# MembershipsLevel List Command

List memberships levels' info, such as id, displayName, etc.

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
# List membership levels
yutu membershipsLevel list
```
