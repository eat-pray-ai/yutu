package yutuber

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideo     error = errors.New("failed to get video")
	errInsertVideo  error = errors.New("failed to insert video")
	errUpdateVideo  error = errors.New("failed to update video")
	errOpenFile     error = errors.New("failed to open file")
	errSetThumbnail error = errors.New("failed to set thumbnail")
)

type Video struct {
	id         string
	path       string
	title      string
	desc       string
	tags       []string
	language   string
	thumbnail  string
	forKids    bool
	restricted bool
	embeddable bool
	category   string
	privacy    string
	channelId  string
}

type VideoService interface {
	List()
	Insert()
	Update()
	get([]string) *youtube.Video
	setThumbnail()
	validate()
}

type VideoOption func(*Video)

func NewVideo(opts ...VideoOption) *Video {
	v := &Video{}
	for _, opt := range opts {
		opt(v)
	}

	v.validate()
	return v
}

func (v *Video) get(parts []string) *youtube.Video {
	service := auth.NewY2BService(youtube.YoutubeReadonlyScope)
	call := service.Videos.List(parts)
	if v.id != "" {
		call = call.Id(v.id)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideo, err), v.id)
	}

	return response.Items[0]
}

func (v *Video) List() {
	video := v.get([]string{"id", "snippet", "status", "statistics"})
	fmt.Printf("          ID: %s\n", video.Id)
	fmt.Printf("       Title: %s\n", video.Snippet.Title)
	fmt.Printf(" Description: %s\n", video.Snippet.Description)
	fmt.Printf("        Tags: %s\n", strings.Join(video.Snippet.Tags, ","))
	fmt.Printf("    language: %s\n", video.Snippet.DefaultLanguage)
	fmt.Printf("  Channel ID: %s\n", video.Snippet.ChannelId)
	fmt.Printf("    Category: %s\n", video.Snippet.CategoryId)
	fmt.Printf("Published At: %s\n", video.Snippet.PublishedAt)
	fmt.Printf("    Privacy: %s\n", video.Status.PrivacyStatus)
	fmt.Printf("   For Kids: %t\n", video.Status.MadeForKids)
	fmt.Printf(" Embeddable: %t\n\n", video.Status.Embeddable)
	fmt.Printf(" Comment Count: %d\n", video.Statistics.CommentCount)
	fmt.Printf(" Dislike Count: %d\n", video.Statistics.DislikeCount)
	fmt.Printf("    Like Count: %d\n", video.Statistics.LikeCount)
	fmt.Printf("Favorite Count: %d\n", video.Statistics.FavoriteCount)
}

func (v *Video) Insert() {
	file, err := os.Open(v.path)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err), v.path)
	}
	defer file.Close()

	upload := &youtube.Video{
		AgeGating: &youtube.VideoAgeGating{
			Restricted: v.restricted,
		},
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

	service := auth.NewY2BService(youtube.YoutubeUploadScope)
	call := service.Videos.Insert([]string{"agegating,snippet,status"}, upload)

	video, err := call.Media(file).Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertVideo, err))
	}

	if v.thumbnail != "" {
		v.setThumbnail(v.thumbnail, service)
	}

	fmt.Printf("Upload successful! Video ID: %v\n", video.Id)
}

func (v *Video) Update() {
	video := v.get([]string{"id", "snippet", "status"})
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

	service := auth.NewY2BService(youtube.YoutubeScope)
	call := service.Videos.Update([]string{"snippet,status"}, video)
	_, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateVideo, err), v.id)
	}

	if v.thumbnail != "" {
		v.setThumbnail(v.thumbnail, service)
	}
	fmt.Println("Video updated:")
	v.List()
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

func (v *Video) validate() {
	if v.forKids && v.restricted {
		log.Fatalln("Video cannot be both for kids and restricted")
	}
}

func WithVideoId(id string) VideoOption {
	return func(v *Video) {
		v.id = id
	}
}

func WithVideoPath(path string) VideoOption {
	return func(v *Video) {
		v.path = path
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

func WithVideoTags(tags string) VideoOption {
	return func(v *Video) {
		v.tags = strings.Split(tags, ",")
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

func WithVideoForKids(forKids bool) VideoOption {
	return func(v *Video) {
		v.forKids = forKids
	}
}

func WithVideoRestricted(restricted bool) VideoOption {
	return func(v *Video) {
		v.restricted = restricted
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
