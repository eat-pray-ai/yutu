# PlaylistImage Update Command

Update a playlist image. Use this tool when you need to update a playlist image.

## Usage

```bash
yutu playlistImage update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--height` | `-H` | The image height |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` | ID of the content owner, for YouTube content partners |
| `--onBehalfOfContentOwnerChannel` | `-B` | YouTube channel ID linked to the content owner |
| `--output` | `-o` | json\|yaml\|silent |
| `--playlistId` | `-p` | ID of the playlist this image is associated with |
| `--type` | `-t` | The image type (e.g., 'hero') |
| `--width` | `-W` | The image width |

## Examples

```bash
yutu playlistImage update --playlistId PLxxx --type hero --width 2048 --height 1152
```
