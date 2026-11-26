// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetSubscription    = errors.New("failed to get subscription")
	errDeleteSubscription = errors.New("failed to delete subscription")
	errInsertSubscription = errors.New("failed to insert subscription")
)

type Subscription struct {
	common.Fields
	SubscriberChannelId string `yaml:"subscriber_channel_id" json:"subscriber_channel_id,omitempty"`
	Description         string `yaml:"description" json:"description,omitempty"`
	ForChannelId        string `yaml:"for_channel_id" json:"for_channel_id,omitempty"`
	For                 string `yaml:"for" json:"for,omitempty"`
	Order               string `yaml:"order" json:"order,omitempty"`
	Title               string `yaml:"title" json:"title,omitempty"`

	OBOCOC string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type ISubscription[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Insert(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*Subscription)

func NewSubscription(opts ...Option) ISubscription[youtube.Subscription] {
	s := &Subscription{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Subscription) Get() ([]*youtube.Subscription, error) {
	if err := s.EnsureService(); err != nil {
		return nil, err
	}
	call := s.Service.Subscriptions.List(s.Parts)
	if len(s.Ids) > 0 {
		call = call.Id(s.Ids...)
	}
	if s.ChannelId != "" {
		call = call.ChannelId(s.ChannelId)
	}
	if s.ForChannelId != "" {
		call = call.ForChannelId(s.ForChannelId)
	}
	switch s.For {
	case "mine":
		call = call.Mine(true)
	case "myRecentSubscribers":
		call = call.MyRecentSubscribers(true)
	case "mySubscribers":
		call = call.MySubscribers(true)
	}
	if s.OBOCO != "" {
		call = call.OnBehalfOfContentOwner(s.OBOCO)
	}
	if s.OBOCOC != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.OBOCOC)
	}
	if s.Order != "" {
		call = call.Order(s.Order)
	}

	return common.Paginate(
		&s.Fields, call,
		func(r *youtube.SubscriptionListResponse) ([]*youtube.Subscription, string) {
			return r.Items, r.NextPageToken
		}, errGetSubscription,
	)
}

func (s *Subscription) List(writer io.Writer) error {
	subscriptions, err := s.Get()
	if err != nil && subscriptions == nil {
		return err
	}

	return errors.Join(
		err,
		common.PrintList(
			s.Output, subscriptions, writer,
			table.Row{"ID", "Kind", "Resource ID", "Channel Title"},
			func(sub *youtube.Subscription) table.Row {
				var resourceId string
				switch sub.Snippet.ResourceId.Kind {
				case "youtube#video":
					resourceId = sub.Snippet.ResourceId.VideoId
				case "youtube#channel":
					resourceId = sub.Snippet.ResourceId.ChannelId
				case "youtube#playlist":
					resourceId = sub.Snippet.ResourceId.PlaylistId
				}
				return table.Row{sub.Id, sub.Snippet.ResourceId.Kind, resourceId, sub.Snippet.Title}
			},
		),
	)
}

func (s *Subscription) Insert(writer io.Writer) error {
	if err := s.EnsureService(); err != nil {
		return err
	}
	subscription := &youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ChannelId:   s.SubscriberChannelId,
			Description: s.Description,
			ResourceId: &youtube.ResourceId{
				ChannelId: s.ChannelId,
			},
			Title: s.Title,
		},
	}

	call := s.Service.Subscriptions.Insert([]string{"snippet"}, subscription)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertSubscription, err)
	}

	return common.PrintResult(
		s.Output, res, writer, "Subscription inserted: %s\n", res.Id,
	)
}

func (s *Subscription) Delete(writer io.Writer) (errs error) {
	if errs = s.EnsureService(); errs != nil {
		return errs
	}
	for _, id := range s.Ids {
		call := s.Service.Subscriptions.Delete(id)
		err := call.Do()
		if err != nil {
			errs = errors.Join(errs, errDeleteSubscription, err)
			continue
		}

		_, err = fmt.Fprintf(writer, "Subscription %s deleted", id)
		errs = errors.Join(errs, err)
	}
	return errs
}

func WithSubscriberChannelId(id string) Option {
	return func(s *Subscription) {
		s.SubscriberChannelId = id
	}
}

func WithDescription(description string) Option {
	return func(s *Subscription) {
		s.Description = description
	}
}

func WithForChannelId(forChannelId string) Option {
	return func(s *Subscription) {
		s.ForChannelId = forChannelId
	}
}

func WithFor(f string) Option {
	return func(s *Subscription) {
		s.For = f
	}
}

func WithOBOCOC(obococ string) Option {
	return func(s *Subscription) {
		s.OBOCOC = obococ
	}
}

func WithOrder(order string) Option {
	return func(s *Subscription) {
		s.Order = order
	}
}

func WithTitle(title string) Option {
	return func(s *Subscription) {
		s.Title = title
	}
}

var (
	WithParts      = common.WithParts[*Subscription]
	WithOutput     = common.WithOutput[*Subscription]
	WithService    = common.WithService[*Subscription]
	WithIds        = common.WithIds[*Subscription]
	WithMaxResults = common.WithMaxResults[*Subscription]
	WithChannelId  = common.WithChannelId[*Subscription]

	WithOBOCO = common.WithOBOCO[*Subscription]
)
