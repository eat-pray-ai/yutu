# Caption Insert

Insert a caption. Use this skill to insert a caption to a video.

## Usage

```bash
yutu caption insert [flags]
```

## Flags

| Flag | Shorthand | Required | Description |
|------|-----------|----------|-------------|
| `--audioTrackType` | `-a` |  | unknown\|primary\|commentary\|descriptive (default "unknown") |
| `--file` | `-f` | Yes | Path to save the caption file |
| `--isAutoSynced` | `-A` |  | Whether YouTube synchronized the caption track to the audio track in the video (default true) |
| `--isCC` | `-C` |  | Whether the track contains closed captions for the deaf and hard of hearing |
| `--isDraft` | `-D` |  | whether the caption track is a draft |
| `--isEasyReader` | `-E` |  | Whether caption track is formatted for 'easy reader' |
| `--isLarge` | `-L` |  | Whether the caption track uses large text for the vision-impaired |
| `--jsonPath` | `-j` |  | JSONPath expression to filter the output |
| `--language` | `-l` |  | Language of the caption track |
| `--name` | `-n` |  | Name of the caption track |
| `--onBehalfOf` | `-b` |  | ID of the YouTube account that the content owner is acting on behalf of |
| `--onBehalfOfContentOwner` | `-B` |  | ID of the content owner, for YouTube content partners |
| `--output` | `-o` |  | json\|yaml\|silent |
| `--trackKind` | `-t` |  | standard\|ASR\|forced (default "standard") |
| `--videoId` | `-v` | Yes | ID of the video |

## Examples

```bash
# Insert a caption to a video
yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ
# Insert a caption with language and name
yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ --language en --name English
# Insert a standard caption and publish immediately
yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ --trackKind standard --isDraft=false
```
