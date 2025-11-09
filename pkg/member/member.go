// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

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
	service      *youtube.Service
	errGetMember = errors.New("failed to get member")
)

type member struct {
	MemberChannelId  string `yaml:"member_channel_id" json:"member_channel_id"`
	HasAccessToLevel string `yaml:"has_access_to_level" json:"has_access_to_level"`
	MaxResults       int64  `yaml:"max_results" json:"max_results"`
	Mode             string `yaml:"mode" json:"mode"`
}

type Member[T any] interface {
	List([]string, string, string, io.Writer) error
	Get([]string) ([]*T, error)
}

type Option func(*member)

func NewMember(opts ...Option) Member[youtube.Member] {
	m := &member{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *member) Get(parts []string) ([]*youtube.Member, error) {
	call := service.Members.List(parts)
	if m.MemberChannelId != "" {
		call = call.FilterByMemberChannelId(m.MemberChannelId)
	}
	if m.HasAccessToLevel != "" {
		call = call.HasAccessToLevel(m.HasAccessToLevel)
	}
	if m.Mode != "" {
		call = call.Mode(m.Mode)
	}

	var items []*youtube.Member
	pageToken := ""
	for m.MaxResults > 0 {
		call = call.MaxResults(min(m.MaxResults, pkg.PerPage))
		m.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetMember, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (m *member) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	members, err := m.Get(parts)
	if err != nil && members == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(members, jpath, writer)
	case "yaml":
		utils.PrintYAML(members, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"Channel ID", "Display Name"})
		for _, member := range members {
			tb.AppendRow(
				table.Row{
					member.Snippet.MemberDetails.ChannelId,
					member.Snippet.MemberDetails.DisplayName,
				},
			)
		}
	}
	return err
}

func WithMemberChannelId(channelId string) Option {
	return func(m *member) {
		m.MemberChannelId = channelId
	}
}

func WithHasAccessToLevel(level string) Option {
	return func(m *member) {
		m.HasAccessToLevel = level
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(m *member) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		m.MaxResults = maxResults
	}
}

func WithMode(mode string) Option {
	return func(m *member) {
		m.Mode = mode
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *member) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
