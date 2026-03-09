# Playlist Image Update

Update a playlist image. Use this skill to update a playlist image.

## Usage

```bash
yutu playlistImage update [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--height` | `-H` |  | The image height |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` |  | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` |  | YouTube channel ID linked to the content owner |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--playlistId` | `-p` | Yes | ID of the playlist this image is associated with |
| `--type` | `-t` |  | The image type (e.g., 'hero') |
| `--width` | `-W` |  | The image width |

## Examples

```bash
# Update a playlist image
yutu playlistImage update --playlistId PLxxx --type hero --width 2048 --height 1152
```
