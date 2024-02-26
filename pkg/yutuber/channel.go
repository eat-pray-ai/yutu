package yutuber

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannel    error = errors.New("failed to get channel")
	errUpdateChannel error = errors.New("failed to update channel")
)

type Channel struct {
	id    string
	title string
	desc  string
	user  string
}

var part = []string{"snippet", "statistics"}

type ChannelService interface {
	List()
	Update()
	get() *youtube.Channel
}

type ChannelOption func(*Channel)

func NewChannel(opts ...ChannelOption) *Channel {
	c := &Channel{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Channel) List() {
	channels := c.get()
	for _, channel := range channels {
		fmt.Printf("          ID: %s\n", channel.Id)
		fmt.Printf("       Title: %s\n", channel.Snippet.Title)
		fmt.Printf(" Description: %s\n", channel.Snippet.Description)
		fmt.Printf("Published At: %s\n", channel.Snippet.PublishedAt)
		fmt.Printf("     Country: %s\n\n", channel.Snippet.Country)
		fmt.Printf("      View Count: %d\n", channel.Statistics.ViewCount)
		fmt.Printf("Subscriber Count: %d\n", channel.Statistics.SubscriberCount)
		fmt.Printf("     Video Count: %d\n\n", channel.Statistics.VideoCount)
	}
}

func (c *Channel) Update() {
	channel := c.get()[0]
	// TODO: is there a better way to check and update?
	if c.title != "" {
		channel.Snippet.Title = c.title
	}
	if c.desc != "" {
		channel.Snippet.Description = c.desc
	}

	service := auth.NewY2BService(youtube.YoutubeScope)
	call := service.Channels.Update(part, channel)
	_, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateChannel, err), c.id)
	}
	fmt.Println("Channel updated:")
	c.List()
}

func (c *Channel) get() []*youtube.Channel {
	service := auth.NewY2BService(youtube.YoutubeReadonlyScope)
	call := service.Channels.List(part)
	switch {
	case c.id != "":
		call = call.Id(c.id)
	case c.user != "":
		call = call.ForUsername(c.user)
	default:
		call = call.Mine(true)
	}
	resp, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetChannel, err), c.id)
	}

	return resp.Items
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

func WithChannelUser(user string) ChannelOption {
	return func(c *Channel) {
		c.user = user
	}
}
