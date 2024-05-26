package yutuber

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideo     = errors.New("failed to get video")
	errInsertVideo  = errors.New("failed to insert video")
	errUpdateVideo  = errors.New("failed to update video")
	errOpenFile     = errors.New("failed to open file")
	errSetThumbnail = errors.New("failed to set thumbnail")
	errRating       = errors.New("failed to rate video")
	errGetRating    = errors.New("failed to get rating")
)

type video struct {
	id          string
	autoLevels  string
	file        string
	title       string
	description string
	hl          string
	tags        []string
	language    string
	locale      string
	license     string
	thumbnail   string
	rating      string
	chart       string
	channelId   string
	playlistId  string
	categoryId  string
	privacy     string
	forKids     bool
	embeddable  bool
	publishAt   string
	regionCode  string
	stabilize   string
	maxHeight   int64
	maxWidth    int64
	maxResults  int64

	notifySubscribers             bool
	publicStatsViewable           bool
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
}

type Video interface {
	List([]string, string)
	Insert()
	Update()
	Rate()
	GetRating()
	get([]string) []*youtube.Video
	setThumbnail(string, *youtube.Service)
}

type VideoOption func(*video)

func NewVideo(opts ...VideoOption) Video {
	v := &video{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *video) get(parts []string) []*youtube.Video {
	call := service.Videos.List(parts)
	if v.id != "" {
		call = call.Id(v.id)
	}
	if v.chart != "" {
		call = call.Chart(v.chart)
	}
	if v.hl != "" {
		call = call.Hl(v.hl)
	}
	if v.locale != "" {
		call = call.Locale(v.locale)
	}
	if v.categoryId != "" {
		call = call.VideoCategoryId(v.categoryId)
	}
	if v.rating != "" {
		call = call.MyRating(v.rating)
	}
	if v.regionCode != "" {
		call = call.RegionCode(v.regionCode)
	}
	if v.maxHeight != 0 {
		call = call.MaxHeight(v.maxHeight)
	}
	if v.maxWidth != 0 {
		call = call.MaxWidth(v.maxWidth)
	}
	if v.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.onBehalfOfContentOwner)
	}
	call = call.MaxResults(v.maxResults)

	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideo, err), v.id)
	}

	return response.Items
}

func (v *video) List(parts []string, output string) {
	videos := v.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(videos)
	case "yaml":
		utils.PrintYAML(videos)
	default:
		fmt.Println("ID\tTitle")
		for _, video := range videos {
			fmt.Printf("%s\t%s\n", video.Id, video.Snippet.Title)
		}
	}
}

func (v *video) Insert() {
	file, err := os.Open(v.file)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err), v.file)
	}
	defer file.Close()

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                v.title,
			Description:          v.description,
			Tags:                 v.tags,
			CategoryId:           v.categoryId,
			ChannelId:            v.channelId,
			DefaultLanguage:      v.language,
			DefaultAudioLanguage: v.language,
		},
		Status: &youtube.VideoStatus{
			Embeddable:              v.embeddable,
			License:                 v.license,
			SelfDeclaredMadeForKids: v.forKids,
			PublishAt:               v.publishAt,
			PrivacyStatus:           v.privacy,
			PublicStatsViewable:     v.publicStatsViewable,
			ForceSendFields:         []string{"SelfDeclaredMadeForKids"},
		},
	}

	call := service.Videos.Insert([]string{"snippet,status"}, video)

	if v.autoLevels == "true" {
		call = call.AutoLevels(true)
	} else if v.autoLevels == "false" {
		call = call.AutoLevels(false)
	}
	call = call.NotifySubscribers(v.notifySubscribers)
	if v.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.onBehalfOfContentOwner)
	}
	if v.onBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(v.onBehalfOfContentOwnerChannel)
	}
	if v.stabilize == "true" {
		call = call.Stabilize(true)
	} else if v.stabilize == "false" {
		call = call.Stabilize(false)
	}

	res, err := call.Media(file).Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertVideo, err))
	}

	if v.thumbnail != "" {
		v.setThumbnail(v.thumbnail, service)
	}

	if v.playlistId != "" {
		pi := NewPlaylistItem(
			WithPlaylistItemTitle(v.title),
			WithPlaylistItemDescription(v.description),
			WithPlaylistItemVideoId(res.Id),
			WithPlaylistItemPlaylistId(v.playlistId),
			WithPlaylistItemChannelId(v.channelId),
			WithPlaylistItemPrivacy(v.privacy),
		)
		pi.Insert()
	}

	fmt.Println("Video inserted:")
	utils.PrintYAML(res)
}

func (v *video) Update() {
	video := v.get([]string{"id", "snippet", "status"})[0]
	if v.title != "" {
		video.Snippet.Title = v.title
	}
	if v.description != "" {
		video.Snippet.Description = v.description
	}
	if v.tags != nil {
		video.Snippet.Tags = v.tags
	}
	if v.language != "" {
		video.Snippet.DefaultLanguage = v.language
		video.Snippet.DefaultAudioLanguage = v.language
	}
	if v.license != "" {
		video.Status.License = v.license
	}
	if v.categoryId != "" {
		video.Snippet.CategoryId = v.categoryId
	}
	if v.privacy != "" {
		video.Status.PrivacyStatus = v.privacy
	}
	video.Status.Embeddable = v.embeddable

	call := service.Videos.Update([]string{"snippet,status"}, video)
	if v.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateVideo, err), v.id)
	}

	if v.thumbnail != "" {
		v.setThumbnail(v.thumbnail, service)
	}

	if v.playlistId != "" {
		pi := NewPlaylistItem(
			WithPlaylistItemTitle(v.title),
			WithPlaylistItemDescription(v.description),
			WithPlaylistItemVideoId(res.Id),
			WithPlaylistItemPlaylistId(v.playlistId),
			WithPlaylistItemChannelId(v.channelId),
			WithPlaylistItemPrivacy(v.privacy),
		)
		pi.Insert()
	}

	fmt.Println("Video updated:")
	utils.PrintYAML(res)
}

func (v *video) Rate() {
	call := service.Videos.Rate(v.id, v.rating)
	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errRating, err), v.id)
	}
	fmt.Printf("Video %s rated %s\n", v.id, v.rating)
}

func (v *video) GetRating() {
	call := service.Videos.GetRating([]string{v.id})
	if v.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.onBehalfOfContentOwnerChannel)
	}
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetRating, err), v.id)
	}

	utils.PrintYAML(res)
}

func (v *video) setThumbnail(thumbnail string, service *youtube.Service) {
	file, err := os.Open(thumbnail)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err), thumbnail)
	}
	call := service.Thumbnails.Set(v.id).Media(file)
	_, err = call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errSetThumbnail, err))
	}
}

func WithVideoId(id string) VideoOption {
	return func(v *video) {
		v.id = id
	}
}

func WithVideoAutoLevels(autoLevels string) VideoOption {
	return func(v *video) {
		v.autoLevels = autoLevels
	}
}

func WithVideoFile(file string) VideoOption {
	return func(v *video) {
		v.file = file
	}
}

func WithVideoTitle(title string) VideoOption {
	return func(v *video) {
		v.title = title
	}
}

func WithVideoDescription(description string) VideoOption {
	return func(v *video) {
		v.description = description
	}
}

func WithVideoHl(hl string) VideoOption {
	return func(v *video) {
		v.hl = hl
	}
}

func WithVideoTags(tags []string) VideoOption {
	return func(v *video) {
		v.tags = tags
	}
}

func WithVideoLanguage(language string) VideoOption {
	return func(v *video) {
		v.language = language
	}
}

func WithVideoLocale(locale string) VideoOption {
	return func(v *video) {
		v.locale = locale
	}
}

func WithVideoLicense(license string) VideoOption {
	return func(v *video) {
		v.license = license
	}
}

func WithVideoThumbnail(thumbnail string) VideoOption {
	return func(v *video) {
		v.thumbnail = thumbnail
	}
}

func WithVideoRating(rating string) VideoOption {
	return func(v *video) {
		v.rating = rating
	}
}

func WithVideoChart(chart string) VideoOption {
	return func(v *video) {
		v.chart = chart
	}
}

func WithVideoForKids(forKids bool) VideoOption {
	return func(v *video) {
		v.forKids = forKids
	}
}

func WithVideoEmbeddable(embeddable bool) VideoOption {
	return func(v *video) {
		v.embeddable = embeddable
	}
}

func WithVideoCategory(categoryId string) VideoOption {
	return func(v *video) {
		v.categoryId = categoryId
	}
}

func WithVideoPrivacy(privacy string) VideoOption {
	return func(v *video) {
		v.privacy = privacy
	}
}

func WithVideoChannelId(channelId string) VideoOption {
	return func(v *video) {
		v.channelId = channelId
	}
}

func WithVideoPlaylistId(playlistId string) VideoOption {
	return func(v *video) {
		v.playlistId = playlistId
	}
}

func WithVideoPublicStatsViewable(publicStatsViewable bool) VideoOption {
	return func(v *video) {
		v.publicStatsViewable = publicStatsViewable
	}
}

func WithVideoPublishAt(publishAt string) VideoOption {
	return func(v *video) {
		v.publishAt = publishAt
	}
}

func WithVideoRegionCode(regionCode string) VideoOption {
	return func(v *video) {
		v.regionCode = regionCode
	}
}

func WithVideoStabilize(stabilize string) VideoOption {
	return func(v *video) {
		v.stabilize = stabilize
	}
}

func WithVideoMaxHeight(maxHeight int64) VideoOption {
	return func(v *video) {
		v.maxHeight = maxHeight
	}
}

func WithVideoMaxWidth(maxWidth int64) VideoOption {
	return func(v *video) {
		v.maxWidth = maxWidth
	}
}

func WithVideoMaxResults(maxResults int64) VideoOption {
	return func(v *video) {
		v.maxResults = maxResults
	}
}

func WithVideoNotifySubscribers(notifySubscribers bool) VideoOption {
	return func(v *video) {
		v.notifySubscribers = notifySubscribers
	}
}

func WithVideoOnBehalfOfContentOwner(onBehalfOfContentOwner string) VideoOption {
	return func(v *video) {
		v.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithVideoOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) VideoOption {
	return func(v *video) {
		v.onBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}
