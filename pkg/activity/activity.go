// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetActivity = errors.New("failed to get activity")
)

type Activity struct {
	*common.Fields
	Home            *bool  `yaml:"home" json:"home,omitempty"`
	Mine            *bool  `yaml:"mine" json:"mine,omitempty"`
	PublishedAfter  string `yaml:"published_after" json:"published_after,omitempty"`
	PublishedBefore string `yaml:"published_before" json:"published_before,omitempty"`
	RegionCode      string `yaml:"region_code" json:"region_code,omitempty"`
}

type IActivity[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*Activity)

func NewActivity(opts ...Option) IActivity[youtube.Activity] {
	a := &Activity{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *Activity) Get() ([]*youtube.Activity, error) {
	if err := a.EnsureService(); err != nil {
		return nil, err
	}
	call := a.Service.Activities.List(a.Parts)
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

	return common.Paginate(
		a.Fields, call,
		func(r *youtube.ActivityListResponse) ([]*youtube.Activity, string) {
			return r.Items, r.NextPageToken
		}, errGetActivity,
	)
}

func (a *Activity) List(writer io.Writer) error {
	activities, err := a.Get()
	if err != nil && activities == nil {
		return err
	}

	common.PrintList(
		a.Output, activities, writer, table.Row{"ID", "Title", "Type", "Time"},
		func(a *youtube.Activity) table.Row {
			return table.Row{a.Id, a.Snippet.Title, a.Snippet.Type, a.Snippet.PublishedAt}
		},
	)
	return err
}

func WithHome(home *bool) Option {
	return func(a *Activity) {
		if home != nil {
			a.Home = home
		}
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

var (
	WithChannelId  = common.WithChannelId[*Activity]
	WithMaxResults = common.WithMaxResults[*Activity]
	WithParts      = common.WithParts[*Activity]
	WithOutput     = common.WithOutput[*Activity]
	WithService    = common.WithService[*Activity]
)
