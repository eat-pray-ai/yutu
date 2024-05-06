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
	errGetActivity = errors.New("failed to get activity")
)

type activity struct {
	channelId       string
	home            string
	maxResults      int64
	mine            string
	publishedAfter  string
	publishedBefore string
	regionCode      string
}

type Activity interface {
	List([]string, string)
	get([]string) []*youtube.Activity
}

type ActivityOption func(*activity)

func NewActivity(opts ...ActivityOption) Activity {
	a := &activity{}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *activity) get(parts []string) []*youtube.Activity {
	service := auth.NewY2BService()
	call := service.Activities.List(parts)
	if a.channelId != "" {
		call = call.ChannelId(a.channelId)
	}

	if a.home == "true" {
		call = call.Home(true)
	} else if a.home == "false" {
		call = call.Home(false)
	}

	call.MaxResults(a.maxResults)

	if a.publishedAfter != "" {
		call.PublishedAfter(a.publishedAfter)
	}

	if a.publishedBefore != "" {
		call.PublishedBefore(a.publishedBefore)
	}

	if a.regionCode != "" {
		call.RegionCode(a.regionCode)
	}

	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetActivity, err))
	}

	return response.Items
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

func WithActivityChannelId(channelId string) ActivityOption {
	return func(a *activity) {
		a.channelId = channelId
	}
}

func WithActivityHome(home string) ActivityOption {
	return func(a *activity) {
		a.home = home
	}
}

func WithActivityMaxResults(maxResults int64) ActivityOption {
	return func(a *activity) {
		a.maxResults = maxResults
	}
}

func WithActivityMine(mine string) ActivityOption {
	return func(a *activity) {
		a.mine = mine
	}
}

func WithActivityPublishedAfter(publishedAfter string) ActivityOption {
	return func(a *activity) {
		a.publishedAfter = publishedAfter
	}
}

func WithActivityPublishedBefore(publishedBefore string) ActivityOption {
	return func(a *activity) {
		a.publishedBefore = publishedBefore
	}
}

func WithActivityRegionCode(regionCode string) ActivityOption {
	return func(a *activity) {
		a.regionCode = regionCode
	}
}
