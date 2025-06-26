package activity

import (
	"errors"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
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
	List([]string, string, io.Writer) error
	Get([]string) ([]*youtube.Activity, error)
}

type Option func(*activity)

func NewActivity(opts ...Option) Activity {
	a := &activity{}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *activity) Get(parts []string) ([]*youtube.Activity, error) {
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
		return nil, errors.Join(errGetActivity, err)
	}

	return res.Items, nil
}

func (a *activity) List(
	parts []string, output string, writer io.Writer,
) error {
	activities, err := a.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(activities, writer)
	case "yaml":
		utils.PrintYAML(activities, writer)
	default:
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Title", "Type"})
		for _, activity := range activities {
			tb.AppendRow(
				table.Row{activity.Id, activity.Snippet.Title, activity.Snippet.Type},
			)
		}
	}
	return nil
}

func WithChannelId(channelId string) Option {
	return func(a *activity) {
		a.ChannelId = channelId
	}
}

func WithHome(home *bool) Option {
	return func(a *activity) {
		if home != nil {
			a.Home = home
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(a *activity) {
		a.MaxResults = maxResults
	}
}

func WithMine(mine *bool) Option {
	return func(a *activity) {
		if mine != nil {
			a.Mine = mine
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
	return func(_ *activity) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
