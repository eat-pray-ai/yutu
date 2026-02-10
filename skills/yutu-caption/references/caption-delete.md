# Caption Delete Command

Delete captions of a video by ids.

## Usage

```bash
yutu caption delete [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--ids` | `-i` | IDs of the captions to delete |
| `--onBehalfOf` | `-b` | |
| `--onBehalfOfContentOwner` | `-B` | |

## Examples

```bash
# Delete a caption
yutu caption delete --ids CAPTION_ID
```
