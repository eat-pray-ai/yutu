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
	errGetVideo     error = errors.New("failed to get video")
	errInsertVideo  error = errors.New("failed to insert video")
	errUpdateVideo  error = errors.New("failed to update video")
	errOpenFile     error = errors.New("failed to open file")
	errSetThumbnail error = errors.New("failed to set thumbnail")
	errGetRating    error = errors.New("failed to get rating")
)

type Video struct {
	id         string
	file       string
	title      string
	desc       string
	tags       []string
	language   string
	thumbnail  string
	myRating   string
	chart      string
	channelId  string
	playlistId string
	category   string
	privacy    string
	forKids    bool
	embeddable bool
}

type VideoService interface {
	List([]string, string)
	Insert()
	Update()
	GetRating()
	get([]string) *youtube.Video
	setThumbnail()
	validate()
}

type VideoOption func(*Video)

func NewVideo(opts ...VideoOption) *Video {
	v := &Video{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *Video) get(parts []string) []*youtube.Video {
	call := service.Videos.List(parts)
	if v.id != "" {
		call = call.Id(v.id)
	} else if v.myRating != "" {
		call = call.MyRating(v.myRating)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideo, err), v.id)
	}

	return response.Items
}

func (v *Video) List(parts []string, output string) {
	videos := v.get([]string{"id", "snippet", "status", "statistics"})
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

func (v *Video) Insert() {
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

func (v *Video) Update() {
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

func (v *Video) GetRating() {
	call := service.Videos.GetRating([]string{v.id})
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetRating, err), v.id)
	}

	utils.PrintJSON(res)
}

func (v *Video) setThumbnail(thumbnail string, service *youtube.Service) {
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
	return func(v *Video) {
		v.id = id
	}
}

func WithVideoFile(file string) VideoOption {
	return func(v *Video) {
		v.file = file
	}
}

func WithVideoTitle(title string) VideoOption {
	return func(v *Video) {
		v.title = title
	}
}

func WithVideoDesc(desc string) VideoOption {
	return func(v *Video) {
		v.desc = desc
	}
}

func WithVideoTags(tags []string) VideoOption {
	return func(v *Video) {
		v.tags = tags
	}
}

func WithVideoLanguage(language string) VideoOption {
	return func(v *Video) {
		v.language = language
	}
}

func WithVideoThumbnail(thumbnail string) VideoOption {
	return func(v *Video) {
		v.thumbnail = thumbnail
	}
}

func WithVideoMyRating(myRating string) VideoOption {
	return func(v *Video) {
		v.myRating = myRating
	}
}

func WithVideoChart(chart string) VideoOption {
	return func(v *Video) {
		v.chart = chart
	}
}

func WithVideoForKids(forKids bool) VideoOption {
	return func(v *Video) {
		v.forKids = forKids
	}
}

func WithVideoEmbeddable(embeddable bool) VideoOption {
	return func(v *Video) {
		v.embeddable = embeddable
	}
}

func WithVideoCategory(category string) VideoOption {
	return func(v *Video) {
		v.category = category
	}
}

func WithVideoPrivacy(privacy string) VideoOption {
	return func(v *Video) {
		v.privacy = privacy
	}
}

func WithVideoChannelId(channelId string) VideoOption {
	return func(v *Video) {
		v.channelId = channelId
	}
}

func WithVideoPlaylistId(playlistId string) VideoOption {
	return func(v *Video) {
		v.playlistId = playlistId
	}
}
