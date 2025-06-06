package videoCategory

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service             *youtube.Service
	errGetVideoCategory = errors.New("failed to get video categoryId")
)

type videoCategory struct {
	IDs        []string `yaml:"ids" json:"ids"`
	Hl         string   `yaml:"hl" json:"hl"`
	RegionCode string   `yaml:"region_code" json:"region_code"`
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
	if len(vc.IDs) > 0 {
		call = call.Id(vc.IDs...)
	}
	if vc.Hl != "" {
		call = call.Hl(vc.Hl)
	}
	if vc.RegionCode != "" {
		call = call.RegionCode(vc.RegionCode)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(vc, nil)
		log.Fatalln(errors.Join(errGetVideoCategory, err))
	}

	return res.Items
}

func (vc *videoCategory) List(parts []string, output string) {
	videoCategories := vc.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(videoCategories, nil)
	case "yaml":
		utils.PrintYAML(videoCategories, nil)
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

func WithIDs(ids []string) Option {
	return func(vc *videoCategory) {
		vc.IDs = ids
	}
}

func WithHl(hl string) Option {
	return func(vc *videoCategory) {
		vc.Hl = hl
	}
}

func WithRegionCode(regionCode string) Option {
	return func(vc *videoCategory) {
		vc.RegionCode = regionCode
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *videoCategory) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
