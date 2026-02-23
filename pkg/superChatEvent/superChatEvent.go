// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"errors"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetSuperChatEvent = errors.New("failed to get super chat event")
)

type SuperChatEvent struct {
	*common.Fields
	Hl         string `yaml:"hl" json:"hl,omitempty"`
	MaxResults int64  `yaml:"max_results" json:"max_results,omitempty"`
}

type ISuperChatEvent[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*SuperChatEvent)

func NewSuperChatEvent(opts ...Option) ISuperChatEvent[youtube.SuperChatEvent] {
	s := &SuperChatEvent{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *SuperChatEvent) Get() ([]*youtube.SuperChatEvent, error) {
	s.EnsureService()
	call := s.Service.SuperChatEvents.List(s.Parts)
	if s.Hl != "" {
		call = call.Hl(s.Hl)
	}

	var items []*youtube.SuperChatEvent
	pageToken := ""
	for s.MaxResults > 0 {
		call = call.MaxResults(min(s.MaxResults, pkg.PerPage))
		s.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetSuperChatEvent, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (s *SuperChatEvent) List(writer io.Writer) error {
	events, err := s.Get()
	if err != nil && events == nil {
		return err
	}

	switch s.Output {
	case "json":
		utils.PrintJSON(events, s.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(events, s.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Amount", "Comment", "Supporter"})
		for _, event := range events {
			tb.AppendRow(
				table.Row{
					event.Id, event.Snippet.DisplayString, event.Snippet.CommentText,
					event.Snippet.SupporterDetails.DisplayName,
				},
			)
		}
	}
	return err
}

func WithHl(hl string) Option {
	return func(s *SuperChatEvent) {
		s.Hl = hl
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *SuperChatEvent) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		s.MaxResults = maxResults
	}
}

var (
	WithParts    = common.WithParts[*SuperChatEvent]
	WithOutput   = common.WithOutput[*SuperChatEvent]
	WithJsonpath = common.WithJsonpath[*SuperChatEvent]
	WithService  = common.WithService[*SuperChatEvent]
)
