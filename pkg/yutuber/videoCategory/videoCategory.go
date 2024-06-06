package videoCategory

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service             *youtube.Service
	errGetVideoCategory = errors.New("failed to get video categoryId")
)

type videoCategory struct {
	id         string
	hl         string
	regionCode string
}

type VideoCategory interface {
	get([]string) []*youtube.VideoCategory
	List([]string, string)
}

type Option func(*videoCategory)

func NewVideoCategory(opt ...Option) VideoCategory {
	vc := &videoCategory{}
	for _, o := range opt {
		o(vc)
	}
	return vc
}

func (vc *videoCategory) get(parts []string) []*youtube.VideoCategory {
	call := service.VideoCategories.List(parts)
	if vc.id != "" {
		call = call.Id(vc.id)
	}
	if vc.hl != "" {
		call = call.Hl(vc.hl)
	}
	if vc.regionCode != "" {
		call = call.RegionCode(vc.regionCode)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideoCategory, err), vc.regionCode)
	}

	return res.Items
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

func WithId(id string) Option {
	return func(vc *videoCategory) {
		vc.id = id
	}
}

func WithHl(hl string) Option {
	return func(vc *videoCategory) {
		vc.hl = hl
	}
}

func WithRegionCode(regionCode string) Option {
	return func(vc *videoCategory) {
		vc.regionCode = regionCode
	}
}

func WithService() Option {
	return func(vc *videoCategory) {
		service = auth.NewY2BService()
	}
}
