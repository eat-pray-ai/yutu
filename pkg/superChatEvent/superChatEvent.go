// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"errors"
	"io"

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
	if err := s.EnsureService(); err != nil {
		return nil, err
	}
	call := s.Service.SuperChatEvents.List(s.Parts)
	if s.Hl != "" {
		call = call.Hl(s.Hl)
	}

	return common.Paginate(s.Fields, call, func(r *youtube.SuperChatEventListResponse) ([]*youtube.SuperChatEvent, string) {
		return r.Items, r.NextPageToken
	}, errGetSuperChatEvent)
}

func (s *SuperChatEvent) List(writer io.Writer) error {
	events, err := s.Get()
	if err != nil && events == nil {
		return err
	}

	switch s.Output {
	case "json":
		utils.PrintJSON(events, writer)
	case "yaml":
		utils.PrintYAML(events, writer)
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

var (
	WithHl         = common.WithHl[*SuperChatEvent]
	WithMaxResults = common.WithMaxResults[*SuperChatEvent]
	WithParts      = common.WithParts[*SuperChatEvent]
	WithOutput     = common.WithOutput[*SuperChatEvent]
	WithService    = common.WithService[*SuperChatEvent]
)
