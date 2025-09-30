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

type Activity[T any] interface {
	List([]string, string, string, io.Writer) error
	Get([]string) ([]*T, error)
}

type Option func(*activity)

func NewActivity(opts ...Option) Activity[youtube.Activity] {
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

func (a *activity) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	activities, err := a.Get(parts)
	if err != nil && activities == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(activities, jpath, writer)
	case "yaml":
		utils.PrintYAML(activities, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
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
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
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
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
