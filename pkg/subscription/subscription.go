package subscription

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service               *youtube.Service
	errGetSubscription    = errors.New("failed to get subscription")
	errDeleteSubscription = errors.New("failed to delete subscription")
	errInsertSubscription = errors.New("failed to insert subscription")
)

type subscription struct {
	IDs                           []string `yaml:"ids" json:"ids"`
	SubscriberChannelId           string   `yaml:"subscriber_channel_id" json:"subscriber_channel_id"`
	Description                   string   `yaml:"description" json:"description"`
	ChannelId                     string   `yaml:"channel_id" json:"channel_id"`
	ForChannelId                  string   `yaml:"for_channel_id" json:"for_channel_id"`
	MaxResults                    int64    `yaml:"max_results" json:"max_results"`
	Mine                          *bool    `yaml:"mine" json:"mine"`
	MyRecentSubscribers           *bool    `yaml:"my_recent_subscribers" json:"my_recent_subscribers"`
	MySubscribers                 *bool    `yaml:"my_subscribers" json:"my_subscribers"`
	OnBehalfOfContentOwner        string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string   `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
	Order                         string   `yaml:"order" json:"order"`
	Title                         string   `yaml:"title" json:"title"`
}

type Subscription interface {
	Get([]string) ([]*youtube.Subscription, error)
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Delete(io.Writer) error
}

type Option func(*subscription)

func NewSubscription(opts ...Option) Subscription {
	s := &subscription{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *subscription) Get(parts []string) ([]*youtube.Subscription, error) {
	call := service.Subscriptions.List(parts)
	if len(s.IDs) > 0 {
		call = call.Id(s.IDs...)
	}
	if s.ChannelId != "" {
		call = call.ChannelId(s.ChannelId)
	}
	if s.ForChannelId != "" {
		call = call.ForChannelId(s.ForChannelId)
	}
	if s.MaxResults <= 0 {
		s.MaxResults = 1
	}
	call = call.MaxResults(s.MaxResults)

	if s.Mine != nil {
		call = call.Mine(*s.Mine)
	}
	if s.MyRecentSubscribers != nil {
		call = call.MyRecentSubscribers(*s.MyRecentSubscribers)
	}
	if s.MySubscribers != nil {
		call = call.MySubscribers(*s.MySubscribers)
	}

	if s.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}
	if s.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.OnBehalfOfContentOwnerChannel)
	}
	if s.Order != "" {
		call = call.Order(s.Order)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetSubscription, err)
	}

	return res.Items, nil
}

func (s *subscription) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	subscriptions, err := s.Get(parts)
	if err != nil {
		return errors.Join(errGetSubscription, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(subscriptions, jpath, writer)
	case "yaml":
		utils.PrintYAML(subscriptions, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Kind", "Resource ID", "Channel Title"})
		for _, sub := range subscriptions {
			var resourceId string
			switch sub.Snippet.ResourceId.Kind {
			case "youtube#video":
				resourceId = sub.Snippet.ResourceId.VideoId
			case "youtube#channel":
				resourceId = sub.Snippet.ResourceId.ChannelId
			case "youtube#playlist":
				resourceId = sub.Snippet.ResourceId.PlaylistId
			}
			tb.AppendRow(
				table.Row{
					sub.Id, sub.Snippet.ResourceId.Kind, resourceId, sub.Snippet.Title,
				},
			)
		}
	}
	return nil
}

func (s *subscription) Insert(
	output string, jpath string, writer io.Writer,
) error {
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

	call := service.Subscriptions.Insert([]string{"snippet"}, subscription)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertSubscription, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	default:
		_, _ = fmt.Fprintf(writer, "Subscription inserted: %s\n", res.Id)
	}
	return nil
}

func (s *subscription) Delete(writer io.Writer) error {
	for _, id := range s.IDs {
		call := service.Subscriptions.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteSubscription, err)
		}

		_, _ = fmt.Fprintf(writer, "Subscription %s deleted", id)
	}
	return nil
}

func WithIDs(ids []string) Option {
	return func(s *subscription) {
		s.IDs = ids
	}
}

func WithSubscriberChannelId(id string) Option {
	return func(s *subscription) {
		s.SubscriberChannelId = id
	}
}

func WithDescription(description string) Option {
	return func(s *subscription) {
		s.Description = description
	}
}

func WithChannelId(channelId string) Option {
	return func(s *subscription) {
		s.ChannelId = channelId
	}
}

func WithForChannelId(forChannelId string) Option {
	return func(s *subscription) {
		s.ForChannelId = forChannelId
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *subscription) {
		s.MaxResults = maxResults
	}
}

func WithMine(mine *bool) Option {
	return func(s *subscription) {
		if mine != nil {
			s.Mine = mine
		}
	}
}

func WithMyRecentSubscribers(myRecentSubscribers *bool) Option {
	return func(s *subscription) {
		if myRecentSubscribers != nil {
			s.MyRecentSubscribers = myRecentSubscribers
		}
	}
}

func WithMySubscribers(mySubscribers *bool) Option {
	return func(s *subscription) {
		if mySubscribers != nil {
			s.MySubscribers = mySubscribers
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(s *subscription) {
		s.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(s *subscription) {
		s.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithOrder(order string) Option {
	return func(s *subscription) {
		s.Order = order
	}
}

func WithTitle(title string) Option {
	return func(s *subscription) {
		s.Title = title
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *subscription) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
