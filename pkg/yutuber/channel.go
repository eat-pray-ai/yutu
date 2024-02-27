package yutuber

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"

	"google.golang.org/api/youtube/v3"
)

var (
	service          *youtube.Service = auth.NewY2BService()
	part             []string         = []string{"snippet", "statistics"}
	errGetChannel    error            = errors.New("failed to get channel")
	errUpdateChannel error            = errors.New("failed to update channel")
)

type Channel struct {
	id    string
	title string
	desc  string
	user  string
}

type ChannelService interface {
	List([]string, string)
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

func (c *Channel) List(parts []string, output string) {
	channels := c.get()
	switch output {
	case "json":
		utils.PrintJSON(channels)
	case "yaml":
		utils.PrintYAML(channels)
	default:
		fmt.Println("ID\tTitle")
		for _, channel := range channels {
			fmt.Printf("%s\t%s\n", channel.Id, channel.Snippet.Title)
		}
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

	call := service.Channels.Update(part, channel)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateChannel, err), c.id)
	}
	fmt.Println("Channel updated:")
	utils.PrintJSON(res)
}

func (c *Channel) get() []*youtube.Channel {
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
