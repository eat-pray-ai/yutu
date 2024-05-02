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
	service          *youtube.Service
	errGetChannel    error = errors.New("failed to get channel")
	errUpdateChannel error = errors.New("failed to update channel")
)

type channel struct {
	id    string
	title string
	desc  string
	user  string
}

type Channel interface {
	List([]string, string)
	Update()
	get([]string) []*youtube.Channel
}

type channelOption func(*channel)

func NewChannel(opts ...channelOption) Channel {
	c := &channel{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *channel) get(parts []string) []*youtube.Channel {
	call := service.Channels.List(parts)
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

func (c *channel) List(parts []string, output string) {
	channels := c.get(parts)
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

func (c *channel) Update() {
	parts := []string{"snippet"}
	channel := c.get(parts)[0]
	// TODO: is there a better way to check and update?
	if c.title != "" {
		channel.Snippet.Title = c.title
	}
	if c.desc != "" {
		channel.Snippet.Description = c.desc
	}

	call := service.Channels.Update(parts, channel)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateChannel, err), c.id)
	}
	fmt.Println("Channel updated:")
	utils.PrintJSON(res)
}

func WithChannelID(id string) channelOption {
	return func(c *channel) {
		c.id = id
	}
}

func WithChannelTitle(title string) channelOption {
	return func(c *channel) {
		c.title = title
	}
}

func WithChannelDesc(desc string) channelOption {
	return func(c *channel) {
		c.desc = desc
	}
}

func WithChannelUser(user string) channelOption {
	return func(c *channel) {
		c.user = user
	}
}
