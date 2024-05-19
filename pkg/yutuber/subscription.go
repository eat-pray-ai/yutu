package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	errGetSubscription = errors.New("failed to get subscription")
)

type subscription struct {
	id                            string
	channelId                     string
	forChannelId                  string
	maxResult                     int64
	mine                          string
	myRecentSubscribers           string
	mySubscribers                 string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
}

type Subscription interface {
	get([]string) []*youtube.Subscription
	List([]string, string)
	// Insert()
	// Delete()
}

type SubscriptionOption func(*subscription)

func NewSubscription(opts ...SubscriptionOption) Subscription {
	s := &subscription{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *subscription) get(parts []string) []*youtube.Subscription {
	call := service.Subscriptions.List(parts)
	if s.id != "" {
		call = call.Id(s.id)
	}
	if s.channelId != "" {
		call = call.ChannelId(s.channelId)
	}
	if s.forChannelId != "" {
		call = call.ForChannelId(s.forChannelId)
	}
	call = call.MaxResults(s.maxResult)

	if s.mine == "true" {
		call = call.Mine(true)
	} else if s.mine == "false" {
		call = call.Mine(false)
	}
	if s.myRecentSubscribers == "true" {
		call = call.MyRecentSubscribers(true)
	} else if s.myRecentSubscribers == "false" {
		call = call.MyRecentSubscribers(false)
	}
	if s.mySubscribers == "true" {
		call = call.MySubscribers(true)
	} else if s.mySubscribers == "false" {
		call = call.MySubscribers(false)
	}

	if s.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.onBehalfOfContentOwner)
	}
	if s.onBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.onBehalfOfContentOwnerChannel)
	}
	if s.order != "" {
		call = call.Order(s.order)
	}

	res, err := call.Do()
	if err != nil {
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
		fmt.Println("Channel ID\tChannel Title")
		for _, subscription := range subscriptions {
			fmt.Printf("%s\t%s\n", subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.Title)
		}
	}
}

func WithSubscriptionId(id string) SubscriptionOption {
	return func(s *subscription) {
		s.id = id
	}
}

func WithSubscriptionChannelId(channelId string) SubscriptionOption {
	return func(s *subscription) {
		s.channelId = channelId
	}
}

func WithSubscriptionForChannelId(forChannelId string) SubscriptionOption {
	return func(s *subscription) {
		s.forChannelId = forChannelId
	}
}

func WithSubscriptionMaxResult(maxResult int64) SubscriptionOption {
	return func(s *subscription) {
		s.maxResult = maxResult
	}
}

func WithSubscriptionMine(mine string) SubscriptionOption {
	return func(s *subscription) {
		s.mine = mine
	}
}

func WithSubscriptionMyRecentSubscribers(myRecentSubscribers string) SubscriptionOption {
	return func(s *subscription) {
		s.myRecentSubscribers = myRecentSubscribers
	}
}

func WithSubscriptionMySubscribers(mySubscribers string) SubscriptionOption {
	return func(s *subscription) {
		s.mySubscribers = mySubscribers
	}
}

func WithSubscriptionOnBehalfOfContentOwner(onBehalfOfContentOwner string) SubscriptionOption {
	return func(s *subscription) {
		s.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithSubscriptionOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) SubscriptionOption {
	return func(s *subscription) {
		s.onBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}
