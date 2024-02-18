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
	errInsertVideo error = errors.New("failed to insert video")
	errOpenVideo   error = errors.New("failed to open video")
)

type Video struct {
	id         string
	path       string
	title      string
	desc       string
	tags       []string
	language   string
	forKids    bool
	restricted bool
	embeddable bool
	category   string
	privacy    string
	channelId  string
	service    *youtube.Service
}

type VideoService interface {
	Insert()
	validate()
}

type VideoOption func(*Video)

func NewVideo(opts ...VideoOption) *Video {
	v := &Video{
		service: auth.NewY2BService(youtube.YoutubeUploadScope),
	}

	for _, opt := range opts {
		opt(v)
	}

	v.validate()
	return v
}

func (v *Video) Insert() {
	file, err := os.Open(v.path)
	if err != nil {
		log.Fatalln(errors.Join(errOpenVideo, err), v.path)
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

	call := v.service.Videos.Insert([]string{"agegating,snippet,status"}, upload)

	response, err := call.Media(file).Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertVideo, err))
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
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
