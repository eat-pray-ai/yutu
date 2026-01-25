// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

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
	errGetSuperChatEvent = errors.New("failed to get super chat event")
)

type SuperChatEvent struct {
	Hl         string   `yaml:"hl" json:"hl"`
	MaxResults int64    `yaml:"max_results" json:"max_results"`
	Parts      []string `yaml:"parts" json:"parts"`
	Output     string   `yaml:"output" json:"output"`
	Jsonpath   string   `yaml:"jsonpath" json:"jsonpath"`
	service    *youtube.Service
}

type ISuperChatEvent[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	preRun()
}

type Option func(*SuperChatEvent)

func NewSuperChatEvent(opts ...Option) ISuperChatEvent[youtube.SuperChatEvent] {
	s := &SuperChatEvent{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *SuperChatEvent) preRun() {
	if s.service == nil {
		s.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (s *SuperChatEvent) Get() ([]*youtube.SuperChatEvent, error) {
	s.preRun()
	call := s.service.SuperChatEvents.List(s.Parts)
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

func WithParts(parts []string) Option {
	return func(s *SuperChatEvent) {
		s.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(s *SuperChatEvent) {
		s.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(s *SuperChatEvent) {
		s.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(s *SuperChatEvent) {
		s.service = svc
	}
}
