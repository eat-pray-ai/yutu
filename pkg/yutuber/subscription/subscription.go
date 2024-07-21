package subscription

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service               *youtube.Service
	errGetSubscription    = errors.New("failed to get subscription")
	errDeleteSubscription = errors.New("failed to delete subscription")
	errInsertSubscription = errors.New("failed to insert subscription")
)

type subscription struct {
	ID                            string `yaml:"id" json:"id"`
	SubscriberChannelId           string `yaml:"subscriber_channel_id" json:"subscriber_channel_id"`
	Description                   string `yaml:"description" json:"description"`
	ChannelId                     string `yaml:"channel_id" json:"channel_id"`
	ForChannelId                  string `yaml:"for_channel_id" json:"for_channel_id"`
	MaxResults                    int64  `yaml:"max_results" json:"max_results"`
	Mine                          *bool  `yaml:"mine" json:"mine"`
	MyRecentSubscribers           *bool  `yaml:"my_recent_subscribers" json:"my_recent_subscribers"`
	MySubscribers                 *bool  `yaml:"my_subscribers" json:"my_subscribers"`
	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
	Order                         string `yaml:"order" json:"order"`
	Title                         string `yaml:"title" json:"title"`
}

type Subscription interface {
	get([]string) []*youtube.Subscription
	List([]string, string)
	Insert(silent bool)
	Delete()
}

type Option func(*subscription)

func NewSubscription(opts ...Option) Subscription {
	s := &subscription{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *subscription) get(parts []string) []*youtube.Subscription {
	call := service.Subscriptions.List(parts)
	if s.ID != "" {
		call = call.Id(s.ID)
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
		utils.PrintJSON(s)
		log.Fatalln(errors.Join(errGetSubscription, err))
	}

	return res.Items
}

func (s *subscription) List(parts []string, output string) {
	subscriptions := s.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(subscriptions)
	case "yaml":
		utils.PrintYAML(subscriptions)
	default:
		fmt.Println("ID\tChannel ID\tChannel Title")
		for _, subscription := range subscriptions {
			fmt.Printf(
				"%s\t%s\t%s\n", subscription.Id,
				subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.Title,
			)
		}
	}
}

func (s *subscription) Insert(silent bool) {
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
		utils.PrintJSON(s)
		log.Fatalln(errors.Join(errInsertSubscription, err))
	}

	if !silent {
		utils.PrintYAML(res)
	}
}

func (s *subscription) Delete() {
	call := service.Subscriptions.Delete(s.ID)
	err := call.Do()
	if err != nil {
		utils.PrintJSON(s)
		log.Fatalln(errors.Join(errDeleteSubscription, err))
	}

	fmt.Printf("Subscription %s deleted", s.ID)
}

func WithID(id string) Option {
	return func(s *subscription) {
		s.ID = id
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

func WithMine(mine bool, changed bool) Option {
	return func(s *subscription) {
		if changed {
			s.Mine = &mine
		}
	}
}

func WithMyRecentSubscribers(myRecentSubscribers bool, changed bool) Option {
	return func(s *subscription) {
		if changed {
			s.MyRecentSubscribers = &myRecentSubscribers
		}
	}
}

func WithMySubscribers(mySubscribers bool, changed bool) Option {
	return func(s *subscription) {
		if changed {
			s.MySubscribers = &mySubscribers
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
	return func(s *subscription) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
