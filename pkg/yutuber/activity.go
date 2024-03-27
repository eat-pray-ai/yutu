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

type Activity struct {
	channelId string
}

type ActivityService interface {
	List([]string, string)
	get() *youtube.Activity
}

type ActivityOption func(*Activity)

func NewActivity(opts ...ActivityOption) *Activity {
	a := &Activity{}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Activity) get() []*youtube.Activity {
	service := auth.NewY2BService()
	call := service.Activities.List(part)
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

func (a *Activity) List(parts []string, output string) {
	activities := a.get()
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
	return func(a *Activity) {
		a.channelId = channelId
	}
}
