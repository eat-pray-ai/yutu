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
	errGetActivity error = errors.New("failed to get activity")
)

type activity struct {
	channelId string
}

type Activity interface {
	List([]string, string)
	get([]string) []*youtube.Activity
}

type activityOption func(*activity)

func NewActivity(opts ...activityOption) Activity {
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
	} else {
		call = call.Mine(true)
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

func WithActivityChannelId(channelId string) activityOption {
	return func(a *activity) {
		a.channelId = channelId
	}
}
