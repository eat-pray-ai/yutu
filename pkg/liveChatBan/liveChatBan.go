// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

var (
	errInsertLiveChatBan = errors.New("failed to insert live chat ban")
	errDeleteLiveChatBan = errors.New("failed to delete live chat ban")
)

type LiveChatBan struct {
	common.Fields
	LiveChatId          string `yaml:"live_chat_id" json:"live_chat_id,omitempty"`
	BannedUserChannelId string `yaml:"banned_user_channel_id" json:"banned_user_channel_id,omitempty"`
	BanDurationSeconds  uint64 `yaml:"ban_duration_seconds" json:"ban_duration_seconds,omitempty"`
	BanType             string `yaml:"ban_type" json:"ban_type,omitempty"`
}

type ILiveChatBan interface {
	Insert(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*LiveChatBan)

func NewLiveChatBan(opts ...Option) ILiveChatBan {
	b := &LiveChatBan{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *LiveChatBan) Insert(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	ban := &youtube.LiveChatBan{
		Snippet: &youtube.LiveChatBanSnippet{
			LiveChatId: b.LiveChatId,
			BannedUserDetails: &youtube.ChannelProfileDetails{
				ChannelId: b.BannedUserChannelId,
			},
			Type:               b.BanType,
			BanDurationSeconds: b.BanDurationSeconds,
		},
	}

	call := b.Service.LiveChatBans.Insert(b.Parts, ban)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertLiveChatBan, err)
	}

	common.PrintResult(
		b.Output, res, writer, "Live chat ban inserted: %s\n", res.Id,
	)
	return nil
}

func (b *LiveChatBan) Delete(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	for _, id := range b.Ids {
		call := b.Service.LiveChatBans.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteLiveChatBan, err)
		}

		_, _ = fmt.Fprintf(writer, "Live chat ban %s deleted\n", id)
	}
	return nil
}

func WithLiveChatId(liveChatId string) Option {
	return func(b *LiveChatBan) {
		b.LiveChatId = liveChatId
	}
}

func WithBannedUserChannelId(channelId string) Option {
	return func(b *LiveChatBan) {
		b.BannedUserChannelId = channelId
	}
}

func WithBanDurationSeconds(duration uint64) Option {
	return func(b *LiveChatBan) {
		b.BanDurationSeconds = duration
	}
}

func WithBanType(banType string) Option {
	return func(b *LiveChatBan) {
		b.BanType = banType
	}
}

var (
	WithParts   = common.WithParts[*LiveChatBan]
	WithOutput  = common.WithOutput[*LiveChatBan]
	WithService = common.WithService[*LiveChatBan]
	WithIds     = common.WithIds[*LiveChatBan]
)
