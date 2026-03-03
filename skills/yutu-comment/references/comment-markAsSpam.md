# Comment MarkAsSpam Command

Mark comments as spam.

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
# Mark a comment as spam
yutu comment markAsSpam --ids abc123
# Mark multiple comments as spam
yutu comment markAsSpam --ids abc123,def456
```
