# PlaylistImage Insert Command

Insert a YouTube playlist image for a given playlist id.

## Usage

```bash
yutu playlistImage insert [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--file` | `-f` | Path to the image file |
| `--height` | `-H` | The image height |
| `--jsonPath` | `-j` | JSONPath expression to filter the output |
| `--onBehalfOfContentOwner` | `-b` | |
| `--onBehalfOfContentOwnerChannel` | `-B` | |
| `--output` | `-o` | json\|yaml\|silent |
| `--playlistId` | `-p` | ID of the playlist this image is associated with |
| `--type` | `-t` | The image type (e.g., 'hero') |
| `--width` | `-W` | The image width |

## Examples

```bash
# Insert a new image for a playlist
yutu playlistImage insert --playlistId PLAYLIST_ID --file image.png
```
