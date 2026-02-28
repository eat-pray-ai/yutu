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
| `--onBehalfOf` | `-b` | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` | ID of the content owner, for YouTube content partners |

## Examples

```bash
# Delete a caption
yutu caption delete --ids CAPTION_ID
```
