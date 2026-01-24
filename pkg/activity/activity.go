// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"errors"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetActivity = errors.New("failed to get activity")
)

type Activity struct {
	service         *youtube.Service
	ChannelId       string `yaml:"channel_id" json:"channel_id"`
	Home            *bool  `yaml:"home" json:"home"`
	MaxResults      int64  `yaml:"max_results" json:"max_results"`
	Mine            *bool  `yaml:"mine" json:"mine"`
	PublishedAfter  string `yaml:"published_after" json:"published_after"`
	PublishedBefore string `yaml:"published_before" json:"published_before"`
	RegionCode      string `yaml:"region_code" json:"region_code"`

	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
}

type IActivity[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	preRun()
}

type Option func(*Activity)

func NewActivity(opts ...Option) IActivity[youtube.Activity] {
	a := &Activity{}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Activity) Get() ([]*youtube.Activity, error) {
	a.preRun()
	call := a.service.Activities.List(a.Parts)
	if a.ChannelId != "" {
		call = call.ChannelId(a.ChannelId)
	}

	if a.Home != nil {
		call = call.Home(*a.Home)
	}

	if a.Mine != nil {
		call = call.Mine(*a.Mine)
	}

	if a.PublishedAfter != "" {
		call = call.PublishedAfter(a.PublishedAfter)
	}

	if a.PublishedBefore != "" {
		call = call.PublishedBefore(a.PublishedBefore)
	}

	if a.RegionCode != "" {
		call = call.RegionCode(a.RegionCode)
	}

	var items []*youtube.Activity
	pageToken := ""
	for a.MaxResults > 0 {
		call = call.MaxResults(min(a.MaxResults, pkg.PerPage))
		a.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetActivity, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (a *Activity) List(writer io.Writer) error {
	activities, err := a.Get()
	if err != nil && activities == nil {
		return err
	}

	switch a.Output {
	case "json":
		utils.PrintJSON(activities, a.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(activities, a.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Title", "Type", "Time"})
		for _, activity := range activities {
			tb.AppendRow(
				table.Row{
					activity.Id, activity.Snippet.Title,
					activity.Snippet.Type, activity.Snippet.PublishedAt,
				},
			)
		}
	}
	return err
}

func (a *Activity) preRun() {
	if a.service == nil {
		a.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func WithChannelId(channelId string) Option {
	return func(a *Activity) {
		a.ChannelId = channelId
	}
}

func WithHome(home *bool) Option {
	return func(a *Activity) {
		if home != nil {
			a.Home = home
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(a *Activity) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		a.MaxResults = maxResults
	}
}

func WithMine(mine *bool) Option {
	return func(a *Activity) {
		if mine != nil {
			a.Mine = mine
		}
	}
}

func WithPublishedAfter(publishedAfter string) Option {
	return func(a *Activity) {
		a.PublishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore string) Option {
	return func(a *Activity) {
		a.PublishedBefore = publishedBefore
	}
}

func WithRegionCode(regionCode string) Option {
	return func(a *Activity) {
		a.RegionCode = regionCode
	}
}

func WithParts(parts []string) Option {
	return func(a *Activity) {
		a.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(a *Activity) {
		a.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(a *Activity) {
		a.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(c *Activity) {
		c.service = svc
	}
}
