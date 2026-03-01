# Comment MarkAsSpam Command

Mark comments as spam. Use this tool when you need to mark comments as spam.

## Usage

```bash
yutu comment markAsSpam [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of comments |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--output` | `-o` | json\|yaml\|silent |

## Examples

```bash
yutu comment markAsSpam --ids abc123
yutu comment markAsSpam --ids abc123,def456
```
