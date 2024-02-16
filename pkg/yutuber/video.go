package yutuber

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/spf13/cobra"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	path     string
	title    string
	desc     string
	category string
	keywords string
	privacy  string
	service  *youtube.Service
}

type VideoService interface {
	Insert()
}

type VideoOption func(*Video)

func NewVideo(opts ...VideoOption) *Video {
	v := &Video{
		service: auth.NewY2BService(youtube.YoutubeUploadScope),
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *Video) Insert() {
	video, err := os.Open(v.path)
	if err != nil {
		log.Fatalf("Error opening %v: %v", v.path, err)
	}
	defer video.Close()

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       v.title,
			Description: v.desc,
			CategoryId:  v.category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: v.privacy},
	}

	if strings.Trim(v.keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(v.keywords, ",")
	}

	call := v.service.Videos.Insert([]string{"snippet,status"}, upload)

	response, err := call.Media(video).Do()
	cobra.CheckErr(err)
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}

func WithPath(path string) VideoOption {
	return func(v *Video) {
		v.path = path
	}
}

func WithTitle(title string) VideoOption {
	return func(v *Video) {
		v.title = title
	}
}

func WithDesc(desc string) VideoOption {
	return func(v *Video) {
		v.desc = desc
	}
}

func WithCategory(category string) VideoOption {
	return func(v *Video) {
		v.category = category
	}
}

func WithKeywords(keywords string) VideoOption {
	return func(v *Video) {
		v.keywords = keywords
	}
}

func WithPrivacy(privacy string) VideoOption {
	return func(v *Video) {
		v.privacy = privacy
	}
}
