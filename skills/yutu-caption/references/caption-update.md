# Caption Update Command

Update a video caption. Use this tool when you need to update a video caption.

## Usage

```bash
yutu caption update [flags]
```

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--audioTrackType` | `-a` | unknown\|primary\|commentary\|descriptive (default "unknown") |
| `--file` | `-f` | Path to save the caption file |
| `--isAutoSynced` | `-A` | Whether YouTube synchronized the caption track to the audio track in the video (default true) |
| `--isCC` | `-C` | Whether the track contains closed captions for the deaf and hard of hearing |
| `--isDraft` | `-D` | whether the caption track is a draft |
| `--isEasyReader` | `-E` | Whether caption track is formatted for 'easy reader' |
| `--isLarge` | `-L` | Whether the caption track uses large text for the vision-impaired |
| `--jsonpath` | `-j` | JSONPath expression to filter the output |
| `--language` | `-l` | Language of the caption track |
| `--name` | `-n` | Name of the caption track |
| `--onBehalfOf` | `-b` | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` | ID of the content owner, for YouTube content partners |
| `--output` | `-o` | json\|yaml\|silent |
| `--trackKind` | `-t` | standard\|ASR\|forced (default "standard") |
| `--videoId` | `-v` | ID of the video |

## Examples

```bash
yutu caption update --videoId dQw4w9WgXcQ --isDraft=false
yutu caption update --videoId dQw4w9WgXcQ --language en --name English
yutu caption update --videoId dQw4w9WgXcQ --file updated.srt
```
