package yutuber

import (
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"

	"google.golang.org/api/youtube/v3"
)

type Channel struct {
	id      string
	title   string
	desc    string
	service *youtube.Service
}

var part = []string{"snippet", "statistics"}

type ChannelService interface {
	List()
	Update()
}

type ChannelOption func(*Channel)

func NewChannel(opts ...ChannelOption) *Channel {
	c := &Channel{
		service: auth.NewY2BService(youtube.YoutubeScope),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Channel) List() {
	channel := c.get()
	fmt.Printf("          ID: %s\n", channel.Id)
	fmt.Printf("       Title: %s\n", channel.Snippet.Title)
	fmt.Printf(" Description: %s\n", channel.Snippet.Description)
	fmt.Printf("Published At: %s\n", channel.Snippet.PublishedAt)
	fmt.Printf("     Country: %s\n\n", channel.Snippet.Country)
	fmt.Printf("      View Count: %d\n", channel.Statistics.ViewCount)
	fmt.Printf("Subscriber Count: %d\n", channel.Statistics.SubscriberCount)
	fmt.Printf("     Video Count: %d\n", channel.Statistics.VideoCount)
}

func (c *Channel) Update() {
	channel := c.get()
	// TODO: is there a better way to check and update?
	if c.title != "" {
		channel.Snippet.Title = c.title
	}
	if c.desc != "" {
		channel.Snippet.Description = c.desc
	}

	call := c.service.Channels.Update(part, channel)
	_, err := call.Do()
	if err != nil {
		log.Fatalf("Failed to update channel: %v", err)
	}

	fmt.Println("Channel updated:")
	c.List()
}

func (c *Channel) get() *youtube.Channel {
	call := c.service.Channels.List(part)
	call = call.Id(c.id)
	resp, err := call.Do()
	if err != nil {
		log.Fatalf("Failed to get channel: %v", err)
	}

	return resp.Items[0]
}

func WithChannelID(id string) ChannelOption {
	return func(c *Channel) {
		c.id = id
	}
}

func WithChannelTitle(title string) ChannelOption {
	return func(c *Channel) {
		c.title = title
	}
}

func WithChannelDesc(desc string) ChannelOption {
	return func(c *Channel) {
		c.desc = desc
	}
}
