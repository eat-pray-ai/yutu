package video

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlistItem"
	"log"
	"os"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service         *youtube.Service
	errGetVideo     = errors.New("failed to get video")
	errInsertVideo  = errors.New("failed to insert video")
	errUpdateVideo  = errors.New("failed to update video")
	errOpenFile     = errors.New("failed to open file")
	errSetThumbnail = errors.New("failed to set thumbnail")
	errRating       = errors.New("failed to rate video")
	errGetRating    = errors.New("failed to get rating")
	errDeleteVideo  = errors.New("failed to delete video")
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
	Delete()
	get([]string) []*youtube.Video
	setThumbnail(string, *youtube.Service)
}

type Option func(*video)

func NewVideo(opts ...Option) Video {
	v := &video{}

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
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(v.title),
			playlistItem.WithDescription(v.description),
			playlistItem.WithVideoId(res.Id),
			playlistItem.WithPlaylistId(v.playlistId),
			playlistItem.WithChannelId(v.channelId),
			playlistItem.WithPrivacy(v.privacy),
			playlistItem.WithService(),
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
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(v.title),
			playlistItem.WithDescription(v.description),
			playlistItem.WithVideoId(res.Id),
			playlistItem.WithPlaylistId(v.playlistId),
			playlistItem.WithChannelId(v.channelId),
			playlistItem.WithPrivacy(v.privacy),
			playlistItem.WithService(),
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

func (v *video) Delete() {
	call := service.Videos.Delete(v.id)
	if v.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeleteVideo, err))
	}
	fmt.Printf("Video %s deleted", v.id)
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

func WithId(id string) Option {
	return func(v *video) {
		v.id = id
	}
}

func WithAutoLevels(autoLevels string) Option {
	return func(v *video) {
		v.autoLevels = autoLevels
	}
}

func WithFile(file string) Option {
	return func(v *video) {
		v.file = file
	}
}

func WithTitle(title string) Option {
	return func(v *video) {
		v.title = title
	}
}

func WithDescription(description string) Option {
	return func(v *video) {
		v.description = description
	}
}

func WithHl(hl string) Option {
	return func(v *video) {
		v.hl = hl
	}
}

func WithTags(tags []string) Option {
	return func(v *video) {
		v.tags = tags
	}
}

func WithLanguage(language string) Option {
	return func(v *video) {
		v.language = language
	}
}

func WithLocale(locale string) Option {
	return func(v *video) {
		v.locale = locale
	}
}

func WithLicense(license string) Option {
	return func(v *video) {
		v.license = license
	}
}

func WithThumbnail(thumbnail string) Option {
	return func(v *video) {
		v.thumbnail = thumbnail
	}
}

func WithRating(rating string) Option {
	return func(v *video) {
		v.rating = rating
	}
}

func WithChart(chart string) Option {
	return func(v *video) {
		v.chart = chart
	}
}

func WithForKids(forKids bool) Option {
	return func(v *video) {
		v.forKids = forKids
	}
}

func WithEmbeddable(embeddable bool) Option {
	return func(v *video) {
		v.embeddable = embeddable
	}
}

func WithCategory(categoryId string) Option {
	return func(v *video) {
		v.categoryId = categoryId
	}
}

func WithPrivacy(privacy string) Option {
	return func(v *video) {
		v.privacy = privacy
	}
}

func WithChannelId(channelId string) Option {
	return func(v *video) {
		v.channelId = channelId
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(v *video) {
		v.playlistId = playlistId
	}
}

func WithPublicStatsViewable(publicStatsViewable bool) Option {
	return func(v *video) {
		v.publicStatsViewable = publicStatsViewable
	}
}

func WithPublishAt(publishAt string) Option {
	return func(v *video) {
		v.publishAt = publishAt
	}
}

func WithRegionCode(regionCode string) Option {
	return func(v *video) {
		v.regionCode = regionCode
	}
}

func WithStabilize(stabilize string) Option {
	return func(v *video) {
		v.stabilize = stabilize
	}
}

func WithMaxHeight(maxHeight int64) Option {
	return func(v *video) {
		v.maxHeight = maxHeight
	}
}

func WithMaxWidth(maxWidth int64) Option {
	return func(v *video) {
		v.maxWidth = maxWidth
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(v *video) {
		v.maxResults = maxResults
	}
}

func WithNotifySubscribers(notifySubscribers bool) Option {
	return func(v *video) {
		v.notifySubscribers = notifySubscribers
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(v *video) {
		v.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(v *video) {
		v.onBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithService() Option {
	return func(v *video) {
		service = auth.NewY2BService()
	}
}
