// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetLiveChatModerator    = errors.New("failed to get live chat moderator")
	errInsertLiveChatModerator = errors.New("failed to insert live chat moderator")
	errDeleteLiveChatModerator = errors.New("failed to delete live chat moderator")
)

type LiveChatModerator struct {
	common.Fields
	LiveChatId         string `yaml:"live_chat_id" json:"live_chat_id,omitempty"`
	ModeratorChannelId string `yaml:"moderator_channel_id" json:"moderator_channel_id,omitempty"`
}

type ILiveChatModerator[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	Insert(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*LiveChatModerator)

func NewLiveChatModerator(opts ...Option) ILiveChatModerator[youtube.LiveChatModerator] {
	m := &LiveChatModerator{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *LiveChatModerator) Get() ([]*youtube.LiveChatModerator, error) {
	if err := m.EnsureService(); err != nil {
		return nil, err
	}
	call := m.Service.LiveChatModerators.List(m.LiveChatId, m.Parts)

	return common.Paginate(
		&m.Fields, call,
		func(r *youtube.LiveChatModeratorListResponse) ([]*youtube.LiveChatModerator, string) {
			return r.Items, r.NextPageToken
		}, errGetLiveChatModerator,
	)
}

func (m *LiveChatModerator) List(writer io.Writer) error {
	moderators, err := m.Get()
	if err != nil && moderators == nil {
		return err
	}

	common.PrintList(
		m.Output, moderators, writer,
		table.Row{"ID", "Channel ID", "Display Name"},
		func(mod *youtube.LiveChatModerator) table.Row {
			return table.Row{
				mod.Id,
				mod.Snippet.ModeratorDetails.ChannelId,
				mod.Snippet.ModeratorDetails.DisplayName,
			}
		},
	)
	return err
}

func (m *LiveChatModerator) Insert(writer io.Writer) error {
	if err := m.EnsureService(); err != nil {
		return err
	}
	moderator := &youtube.LiveChatModerator{
		Snippet: &youtube.LiveChatModeratorSnippet{
			LiveChatId: m.LiveChatId,
			ModeratorDetails: &youtube.ChannelProfileDetails{
				ChannelId: m.ModeratorChannelId,
			},
		},
	}

	call := m.Service.LiveChatModerators.Insert(m.Parts, moderator)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertLiveChatModerator, err)
	}

	common.PrintResult(
		m.Output, res, writer, "Live chat moderator inserted: %s\n", res.Id,
	)
	return nil
}

func (m *LiveChatModerator) Delete(writer io.Writer) error {
	if err := m.EnsureService(); err != nil {
		return err
	}
	for _, id := range m.Ids {
		call := m.Service.LiveChatModerators.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteLiveChatModerator, err)
		}

		_, _ = fmt.Fprintf(writer, "Live chat moderator %s deleted\n", id)
	}
	return nil
}

func WithLiveChatId(liveChatId string) Option {
	return func(m *LiveChatModerator) {
		m.LiveChatId = liveChatId
	}
}

func WithModeratorChannelId(channelId string) Option {
	return func(m *LiveChatModerator) {
		m.ModeratorChannelId = channelId
	}
}

var (
	WithMaxResults = common.WithMaxResults[*LiveChatModerator]
	WithParts      = common.WithParts[*LiveChatModerator]
	WithOutput     = common.WithOutput[*LiveChatModerator]
	WithService    = common.WithService[*LiveChatModerator]
	WithIds        = common.WithIds[*LiveChatModerator]
)
