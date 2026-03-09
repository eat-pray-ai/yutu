# Video Abuse Report Reason List

List video abuse report reasons. Use this skill to list available abuse report reasons.

## Usage

```bash
yutu videoAbuseReportReason list [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--hl` | `-l` |  | Host language |
| `--jsonpath` | `-j` |  | JSONPath expression to filter the output |
| `--output` | `-o` |  | json\|yaml\|table (default "table") |
| `--parts` | `-p` |  | Comma separated parts (default [id,snippet]) |

## Examples

```bash
# List video abuse report reasons
yutu videoAbuseReportReason list --hl en
```
