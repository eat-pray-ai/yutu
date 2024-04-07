package yutuber

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideoCategory error = errors.New("failed to get video category")
)

type VideoCategory struct {
	regionCode string
}

type VideoCategoryService interface {
	get([]string) *youtube.VideoCategory
	List([]string, string)
}

type VideoCategoryOption func(*VideoCategory)

func NewVideoCategory() *VideoCategory {
	service = auth.NewY2BService()
	return &VideoCategory{}
}

func (vc *VideoCategory) get(part []string) []*youtube.VideoCategory {
	call := service.VideoCategories.List(part).RegionCode(vc.regionCode)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideoCategory, err), vc.regionCode)
	}

	return response.Items
}

func (vc *VideoCategory) List(part []string, output string) {
	videoCategories := vc.get(part)
	switch output {
	case "json":
		utils.PrintJSON(videoCategories)
	case "yaml":
		utils.PrintYAML(videoCategories)
	default:
		fmt.Println("ID\tTitle\tAssignable")
		for _, videoCategory := range videoCategories {
			fmt.Printf(
				"%s\t%s\t%t\n", videoCategory.Id,
				videoCategory.Snippet.Title, videoCategory.Snippet.Assignable,
			)
		}
	}
}

func WithRegionCode(regionCode string) VideoCategoryOption {
	return func(vc *VideoCategory) {
		vc.regionCode = regionCode
	}
}
