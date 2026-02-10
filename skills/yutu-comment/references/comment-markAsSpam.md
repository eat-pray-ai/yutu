# Comment MarkAsSpam Command

Mark YouTube comments as spam by ids.

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
# Mark comment as spam
yutu comment markAsSpam --ids COMMENT_ID
```
