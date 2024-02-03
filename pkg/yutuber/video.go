package yutuber

import (
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/util"
)

func VideoInsert(service *youtube.Service, filename string, snippet *youtube.VideoSnippet, keywords string, privacy string) {
	if filename == "" {
		log.Fatal("You must provide a filename of a video file to upload")
	}

	upload := &youtube.Video{
		Snippet: snippet,
		Status:  &youtube.VideoStatus{PrivacyStatus: privacy},
	}

	if strings.Trim(keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(keywords, ",")
	}

	call := service.Videos.Insert([]string{"snippet,status"}, upload)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening %v: %v", filename, err)
	}
	defer file.Close()

	response, err := call.Media(file).Do()
	util.HandleError(err, "")
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}
