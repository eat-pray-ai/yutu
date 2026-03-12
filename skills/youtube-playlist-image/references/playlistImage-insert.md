# Playlist Image Insert

Insert a playlist image. Use this skill to insert a YouTube playlist image for a given playlist ID.

## Usage

```bash
yutu playlistImage insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--file` | `-f` | Yes | Path to the image file |
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
# Insert a playlist cover image
yutu playlistImage insert --file cover.jpg --playlistId PLxxx
# Insert a hero image
yutu playlistImage insert --file cover.png --playlistId PLxxx --type hero
# Insert an image with custom dimensions
yutu playlistImage insert --file cover.jpg --playlistId PLxxx --width 2048 --height 1152
```
