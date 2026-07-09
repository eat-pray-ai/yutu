// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetLiveChatMessage        = errors.New("failed to get live chat message")
	errInsertLiveChatMessage     = errors.New("failed to insert live chat message")
	errDeleteLiveChatMessage     = errors.New("failed to delete live chat message")
	errTransitionLiveChatMessage = errors.New("failed to transition live chat message")
)

type LiveChatMessage struct {
	common.Fields
	LiveChatId  string `yaml:"live_chat_id" json:"live_chat_id,omitempty"`
	MessageText string `yaml:"message_text" json:"message_text,omitempty"`
	Status      string `yaml:"status" json:"status,omitempty"`
}

type ILiveChatMessage[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	Insert(io.Writer) error
	Delete(io.Writer) error
	Transition(io.Writer) error
}

type Option func(*LiveChatMessage)

func NewLiveChatMessage(opts ...Option) ILiveChatMessage[youtube.LiveChatMessage] {
	m := &LiveChatMessage{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *LiveChatMessage) Get() ([]*youtube.LiveChatMessage, error) {
	if err := m.EnsureService(); err != nil {
		return nil, err
	}
	call := m.Service.LiveChatMessages.List(m.LiveChatId, m.Parts)
	if m.Hl != "" {
		call = call.Hl(m.Hl)
	}

	return common.Paginate(
		&m.Fields, call,
		func(r *youtube.LiveChatMessageListResponse) ([]*youtube.LiveChatMessage, string) {
			return r.Items, r.NextPageToken
		}, errGetLiveChatMessage,
	)
}

func (m *LiveChatMessage) List(writer io.Writer) error {
	messages, err := m.Get()
	if err != nil && messages == nil {
		return err
	}

	common.PrintList(
		m.Output, messages, writer,
		table.Row{"ID", "Type", "Author", "Message"},
		func(msg *youtube.LiveChatMessage) table.Row {
			var authorName, msgText string
			if msg.AuthorDetails != nil {
				authorName = msg.AuthorDetails.DisplayName
			}
			if msg.Snippet != nil {
				msgText = msg.Snippet.DisplayMessage
			}
			return table.Row{
				msg.Id,
				msg.Snippet.Type,
				authorName,
				msgText,
			}
		},
	)
	return err
}

func (m *LiveChatMessage) Insert(writer io.Writer) error {
	if err := m.EnsureService(); err != nil {
		return err
	}
	msg := &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			LiveChatId: m.LiveChatId,
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: m.MessageText,
			},
			Type: "textMessageEvent",
		},
	}

	call := m.Service.LiveChatMessages.Insert(m.Parts, msg)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertLiveChatMessage, err)
	}

	common.PrintResult(
		m.Output, res, writer, "Live chat message inserted: %s\n", res.Id,
	)
	return nil
}

func (m *LiveChatMessage) Delete(writer io.Writer) error {
	if err := m.EnsureService(); err != nil {
		return err
	}
	for _, id := range m.Ids {
		call := m.Service.LiveChatMessages.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteLiveChatMessage, err)
		}

		_, _ = fmt.Fprintf(writer, "Live chat message %s deleted\n", id)
	}
	return nil
}

func (m *LiveChatMessage) Transition(writer io.Writer) error {
	if err := m.EnsureService(); err != nil {
		return err
	}
	for _, id := range m.Ids {
		call := m.Service.LiveChatMessages.Transition().Id(id).Status(m.Status)
		res, err := call.Do()
		if err != nil {
			return errors.Join(errTransitionLiveChatMessage, err)
		}

		common.PrintResult(
			m.Output, res, writer,
			"Live chat message %s transitioned to %s\n", id, m.Status,
		)
	}
	return nil
}

func WithLiveChatId(liveChatId string) Option {
	return func(m *LiveChatMessage) {
		m.LiveChatId = liveChatId
	}
}

func WithMessageText(messageText string) Option {
	return func(m *LiveChatMessage) {
		m.MessageText = messageText
	}
}

func WithStatus(status string) Option {
	return func(m *LiveChatMessage) {
		m.Status = status
	}
}

var (
	WithHl         = common.WithHl[*LiveChatMessage]
	WithMaxResults = common.WithMaxResults[*LiveChatMessage]
	WithParts      = common.WithParts[*LiveChatMessage]
	WithOutput     = common.WithOutput[*LiveChatMessage]
	WithService    = common.WithService[*LiveChatMessage]
	WithIds        = common.WithIds[*LiveChatMessage]
)
