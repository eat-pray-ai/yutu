package activity

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service        *youtube.Service
	errGetActivity = errors.New("failed to get activity")
)

type activity struct {
	ChannelId       string `yaml:"channel_id" json:"channel_id"`
	Home            *bool  `yaml:"home" json:"home"`
	MaxResults      int64  `yaml:"max_results" json:"max_results"`
	Mine            *bool  `yaml:"mine" json:"mine"`
	PublishedAfter  string `yaml:"published_after" json:"published_after"`
	PublishedBefore string `yaml:"published_before" json:"published_before"`
	RegionCode      string `yaml:"region_code" json:"region_code"`
}

type Activity interface {
	List([]string, string)
	get([]string) []*youtube.Activity
}

type Option func(*activity)

func NewActivity(opts ...Option) Activity {
	a := &activity{}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *activity) get(parts []string) []*youtube.Activity {
	call := service.Activities.List(parts)
	if a.ChannelId != "" {
		call = call.ChannelId(a.ChannelId)
	}

	if a.Home != nil {
		call = call.Home(*a.Home)
	}

	if a.Mine != nil {
		call = call.Mine(*a.Mine)
	}

	if a.MaxResults <= 0 {
		a.MaxResults = 1
	}
	call.MaxResults(a.MaxResults)

	if a.PublishedAfter != "" {
		call.PublishedAfter(a.PublishedAfter)
	}

	if a.PublishedBefore != "" {
		call.PublishedBefore(a.PublishedBefore)
	}

	if a.RegionCode != "" {
		call.RegionCode(a.RegionCode)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(a)
		log.Fatalln(errors.Join(errGetActivity, err))
	}

	return res.Items
}

func (a *activity) List(parts []string, output string) {
	activities := a.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(activities)
	case "yaml":
		utils.PrintYAML(activities)
	default:
		fmt.Println("ID\tTitle")
		for _, activity := range activities {
			fmt.Printf("%s\t%s\n", activity.Id, activity.Snippet.Title)
		}
	}
}

func WithChannelId(channelId string) Option {
	return func(a *activity) {
		a.ChannelId = channelId
	}
}

func WithHome(home bool, changed bool) Option {
	return func(a *activity) {
		if changed {
			a.Home = &home
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(a *activity) {
		a.MaxResults = maxResults
	}
}

func WithMine(mine bool, changed bool) Option {
	return func(a *activity) {
		if changed {
			a.Mine = &mine
		}
	}
}

func WithPublishedAfter(publishedAfter string) Option {
	return func(a *activity) {
		a.PublishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore string) Option {
	return func(a *activity) {
		a.PublishedBefore = publishedBefore
	}
}

func WithRegionCode(regionCode string) Option {
	return func(a *activity) {
		a.RegionCode = regionCode
	}
}

func WithService(svc *youtube.Service) Option {
	return func(a *activity) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
