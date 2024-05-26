package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	errGetSubscription    = errors.New("failed to get subscription")
	errDeleteSubscription = errors.New("failed to delete subscription")
	errInsertSubscription = errors.New("failed to insert subscription")
)

type subscription struct {
	id                            string
	subscriberChannelId           string
	description                   string
	channelId                     string
	forChannelId                  string
	maxResults                    int64
	mine                          string
	myRecentSubscribers           string
	mySubscribers                 string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	title                         string
}

type Subscription interface {
	get([]string) []*youtube.Subscription
	List([]string, string)
	Insert()
	Delete()
}

type SubscriptionOption func(*subscription)

func NewSubscription(opts ...SubscriptionOption) Subscription {
	service = auth.NewY2BService()
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
	call = call.MaxResults(s.maxResults)

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
		fmt.Println("ID\tChannel ID\tChannel Title")
		for _, subscription := range subscriptions {
			fmt.Printf("%s\t%s\t%s\n", subscription.Id, subscription.Snippet.ResourceId.ChannelId, subscription.Snippet.Title)
		}
	}
}

func (s *subscription) Insert() {
	subscription := &youtube.Subscription{
		Snippet: &youtube.SubscriptionSnippet{
			ChannelId:   s.subscriberChannelId,
			Description: s.description,
			ResourceId: &youtube.ResourceId{
				ChannelId: s.channelId,
			},
			Title: s.title,
		},
	}

	call := service.Subscriptions.Insert([]string{"snippet"}, subscription)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertSubscription, err))
	}
	fmt.Println("Subscription inserted")
	utils.PrintYAML(res)
}

func (s *subscription) Delete() {
	call := service.Subscriptions.Delete(s.id)
	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeleteSubscription, err), s.id)
	}

	fmt.Printf("Subscription %s deleted", s.id)
}

func WithSubscriptionId(id string) SubscriptionOption {
	return func(s *subscription) {
		s.id = id
	}
}

func WithSubscriptionSubscriberChannelId(id string) SubscriptionOption {
	return func(s *subscription) {
		s.subscriberChannelId = id
	}
}

func WithSubscriptionDescription(description string) SubscriptionOption {
	return func(s *subscription) {
		s.description = description
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

func WithSubscriptionMaxResults(maxResults int64) SubscriptionOption {
	return func(s *subscription) {
		s.maxResults = maxResults
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

func WithSubscriptionOrder(order string) SubscriptionOption {
	return func(s *subscription) {
		s.order = order
	}
}

func WithSubscriptionTitle(title string) SubscriptionOption {
	return func(s *subscription) {
		s.title = title
	}
}
