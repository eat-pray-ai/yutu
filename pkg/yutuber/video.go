package yutuber

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/util"
)

type Video struct {
	Path     string
	Title    string
	Desc     string
	Category string
	Keywords string
	Privacy  string
	service  *youtube.Service
}

type VideoService interface {
	Insert()
}

type VideoOption func(*Video)

func NewVideo(ctx context.Context, opts ...VideoOption) *Video {
	authSvc := &authService{Scope: youtube.YoutubeUploadScope}
	service := authSvc.auth(ctx)

	video := &Video{
		service:  service,
	}

	for _, opt := range opts {
		opt(video)
	}

	return video
}

func (v *Video) Insert() {
	video, err := os.Open(v.Path)
	if err != nil {
		log.Fatalf("Error opening %v: %v", v.Path, err)
	}
	defer video.Close()

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       v.Title,
			Description: v.Desc,
			CategoryId:  v.Category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: v.Privacy},
	}

	if strings.Trim(v.Keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(v.Keywords, ",")
	}

	call := v.service.Videos.Insert([]string{"snippet,status"}, upload)

	response, err := call.Media(video).Do()
	util.HandleError(err, "")
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}
