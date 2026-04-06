// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetMember = errors.New("failed to get member")
)

type Member struct {
	*common.Fields
	MemberChannelId  string `yaml:"member_channel_id" json:"member_channel_id,omitempty"`
	HasAccessToLevel string `yaml:"has_access_to_level" json:"has_access_to_level,omitempty"`
	Mode             string `yaml:"mode" json:"mode,omitempty"`
}

type IMember[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*Member)

func NewMember(opts ...Option) IMember[youtube.Member] {
	m := &Member{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Member) Get() ([]*youtube.Member, error) {
	if err := m.EnsureService(); err != nil {
		return nil, err
	}
	call := m.Service.Members.List(m.Parts)
	if m.MemberChannelId != "" {
		call = call.FilterByMemberChannelId(m.MemberChannelId)
	}
	if m.HasAccessToLevel != "" {
		call = call.HasAccessToLevel(m.HasAccessToLevel)
	}
	if m.Mode != "" {
		call = call.Mode(m.Mode)
	}

	return common.Paginate(m.Fields, call, func(r *youtube.MemberListResponse) ([]*youtube.Member, string) {
		return r.Items, r.NextPageToken
	}, errGetMember)
}

func (m *Member) List(writer io.Writer) error {
	members, err := m.Get()
	if err != nil && members == nil {
		return err
	}

	common.PrintList(m.Output, members, writer, table.Row{"Channel ID", "Display Name"}, func(m *youtube.Member) table.Row {
		return table.Row{m.Snippet.MemberDetails.ChannelId, m.Snippet.MemberDetails.DisplayName}
	})
	return err
}

func WithMemberChannelId(channelId string) Option {
	return func(m *Member) {
		m.MemberChannelId = channelId
	}
}

func WithHasAccessToLevel(level string) Option {
	return func(m *Member) {
		m.HasAccessToLevel = level
	}
}

func WithMode(mode string) Option {
	return func(m *Member) {
		m.Mode = mode
	}
}

var (
	WithMaxResults = common.WithMaxResults[*Member]
	WithParts      = common.WithParts[*Member]
	WithOutput     = common.WithOutput[*Member]
	WithService    = common.WithService[*Member]
)
