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

type videoCategory struct {
	regionCode string
}

type VideoCategory interface {
	get([]string) []*youtube.VideoCategory
	List([]string, string)
}

type videoCategoryOption func(*videoCategory)

func NewVideoCategory(opt ...videoCategoryOption) VideoCategory {
	service = auth.NewY2BService()
	vc := &videoCategory{}
	for _, o := range opt {
		o(vc)
	}
	return vc
}

func (vc *videoCategory) get(parts []string) []*youtube.VideoCategory {
	call := service.VideoCategories.List(parts).RegionCode(vc.regionCode)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideoCategory, err), vc.regionCode)
	}

	return response.Items
}

func (vc *videoCategory) List(parts []string, output string) {
	videoCategories := vc.get(parts)
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

func WithRegionCode(regionCode string) videoCategoryOption {
	return func(vc *videoCategory) {
		vc.regionCode = regionCode
	}
}
