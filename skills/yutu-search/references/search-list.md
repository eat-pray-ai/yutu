# Search List Command

Search for YouTube resources.

## Usage

```bash
yutu search list [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--channelId` | Filter on resources belonging to this channelId |
| `--channelType` | channelTypeUnspecified\|any\|show (default "channelTypeUnspecified") |
| `--eventType` | none\|upcoming\|live\|completed (default "none") |
| `--forContentOwner` | Search owned by content owner |
| `--forDeveloper` | Only retrieve videos uploaded using the project id of the authenticated user |
| `--forMine` | Search for the private videos of the authenticated user |
| `--jsonpath` | JSONPath expression to filter the output |
| `--location` | Filter on location of the video |
| `--locationRadius` | Filter on distance from the location |
| `--maxResults` | The maximum number of items that should be returned, 0 for no limit (default 5) |
| `--onBehalfOfContentOwner` | ID of the content owner, for YouTube content partners |
| `--order` | searchSortUnspecified, date, rating, viewCount, relevance, title, videoCount (default "relevance") |
| `--output` | json\|yaml\|table (default "table") |
| `--parts` | Comma separated parts (default [id,snippet]) |
| `--publishedAfter` | Filter on resources published after this date |
| `--publishedBefore` | Filter on resources published before this date |
| `--q` | Textual search terms to match |
| `--regionCode` | Display the content as seen by viewers in this country |
| `--relevanceLanguage` | Return results relevant to this language |
| `--safeSearch` | safeSearchSettingUnspecified\|none\|moderate\|strict (default "moderate") |
| `--topicId` | Restrict results to a particular topic |
| `--types` | Restrict results to a particular set of resource types from One Platform |
| `--videoCaption` | videoCaptionUnspecified\|any\|closedCaption\|none (default "any") |
| `--videoCategoryId` | Filter on videos in a specific category |
| `--videoDefinition` | Filter on the definition of the videos |
| `--videoDimension` | any\|2d\|3d (default "any") |
| `--videoDuration` | videoDurationUnspecified\|any\|short\|medium\|long (default "any") |
| `--videoEmbeddable` | videoEmbeddableUnspecified\|any\|true |
| `--videoLicense` | any\|youtube\|creativeCommon |
| `--videoPaidProductPlacement` | videoPaidProductPlacementUnspecified\|any\|true |
| `--videoSyndicated` | videoSyndicatedUnspecified\|any\|true |
| `--videoType` | videoTypeUnspecified\|any\|movie\|episode |

## Examples

```bash
# Search for videos matching a query
yutu search list --q "Golang tutorial" --types video

# Search for my videos
yutu search list --forMine --types video

# Search for channels
yutu search list --q "Google Developers" --types channel
```
