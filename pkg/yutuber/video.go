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
	id         string
	file       string
	title      string
	desc       string
	tags       []string
	language   string
	thumbnail  string
	rating     string
	chart      string
	channelId  string
	playlistId string
	category   string
	privacy    string
	forKids    bool
	embeddable bool
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
	} else if v.rating != "" {
		call = call.MyRating(v.rating)
	}
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
			Description:          v.desc,
			Tags:                 v.tags,
			CategoryId:           v.category,
			ChannelId:            v.channelId,
			DefaultLanguage:      v.language,
			DefaultAudioLanguage: v.language,
		},
		Status: &youtube.VideoStatus{
			Embeddable:    v.embeddable,
			MadeForKids:   v.forKids,
			PrivacyStatus: v.privacy,
		},
	}

	call := service.Videos.Insert([]string{"snippet,status"}, video)

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
			WithPlaylistItemDesc(v.desc),
			WithPlaylistItemVideoId(res.Id),
			WithPlaylistItemPlaylistId(v.playlistId),
			WithPlaylistItemChannelId(v.channelId),
			WithPlaylistItemPrivacy(v.privacy),
		)
		pi.Insert()
	}

	fmt.Println("Video inserted:")
	utils.PrintJSON(res)
}

func (v *video) Update() {
	video := v.get([]string{"id", "snippet", "status"})[0]
	if v.title != "" {
		video.Snippet.Title = v.title
	}
	if v.desc != "" {
		video.Snippet.Description = v.desc
	}
	if v.tags != nil {
		video.Snippet.Tags = v.tags
	}
	if v.language != "" {
		video.Snippet.DefaultLanguage = v.language
		video.Snippet.DefaultAudioLanguage = v.language
	}
	if v.category != "" {
		video.Snippet.CategoryId = v.category
	}
	if v.privacy != "" {
		video.Status.PrivacyStatus = v.privacy
	}
	video.Status.Embeddable = v.embeddable

	call := service.Videos.Update([]string{"snippet,status"}, video)
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
			WithPlaylistItemDesc(v.desc),
			WithPlaylistItemVideoId(res.Id),
			WithPlaylistItemPlaylistId(v.playlistId),
			WithPlaylistItemChannelId(v.channelId),
			WithPlaylistItemPrivacy(v.privacy),
		)
		pi.Insert()
	}

	fmt.Println("Video updated:")
	utils.PrintJSON(res)
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
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetRating, err), v.id)
	}

	utils.PrintJSON(res)
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

func WithVideoDesc(desc string) VideoOption {
	return func(v *video) {
		v.desc = desc
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

func WithVideoCategory(category string) VideoOption {
	return func(v *video) {
		v.category = category
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
