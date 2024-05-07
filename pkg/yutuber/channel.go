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
	errGetChannel    = errors.New("failed to get channel")
	errUpdateChannel = errors.New("failed to update channel")
)

type channel struct {
	categoryId             string
	forHandle              string
	forUsername            string
	hl                     string
	id                     string
	managedByMe            string
	maxResults             int64
	mine                   string
	mySubscribers          string
	onBehalfOfContentOwner string

	country         string
	customUrl       string
	defaultLanguage string
	description     string
	title           string
}

type Channel interface {
	List([]string, string)
	Update()
	get([]string) []*youtube.Channel
}

type ChannelOption func(*channel)

func NewChannel(opts ...ChannelOption) Channel {
	c := &channel{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *channel) get(parts []string) []*youtube.Channel {
	call := service.Channels.List(parts)
	if c.categoryId != "" {
		call = call.CategoryId(c.categoryId)
	}

	if c.forHandle != "" {
		call = call.ForUsername(c.forHandle)
	}

	if c.forUsername != "" {
		call = call.ForUsername(c.forUsername)
	}

	if c.hl != "" {
		call = call.Hl(c.hl)
	}

	if c.id != "" {
		call = call.Id(c.id)
	}

	if c.managedByMe == "true" {
		call = call.ManagedByMe(true)
	} else if c.managedByMe == "false" {
		call = call.ManagedByMe(false)
	}

	call = call.MaxResults(c.maxResults)

	if c.mine == "true" {
		call = call.Mine(true)
	} else if c.mine == "false" {
		call = call.Mine(false)
	}

	if c.mySubscribers == "true" {
		call = call.MySubscribers(true)
	} else if c.mySubscribers == "false" {
		call = call.MySubscribers(false)
	}

	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
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
	if c.title != "" {
		channel.Snippet.Title = c.title
	}
	if c.description != "" {
		channel.Snippet.Description = c.description
	}

	call := service.Channels.Update(parts, channel)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateChannel, err), c.id)
	}
	fmt.Println("Channel updated:")
	utils.PrintJSON(res)
}

func WithChannelCategoryId(categoryId string) ChannelOption {
	return func(c *channel) {
		c.categoryId = categoryId
	}
}

func WithChannelForHandle(handle string) ChannelOption {
	return func(c *channel) {
		c.forHandle = handle
	}
}

func WithChannelForUsername(username string) ChannelOption {
	return func(c *channel) {
		c.forUsername = username
	}
}

func WithChannelHl(hl string) ChannelOption {
	return func(c *channel) {
		c.hl = hl
	}
}

func WithChannelId(id string) ChannelOption {
	return func(c *channel) {
		c.id = id
	}
}

func WithChannelManagedByMe(managedByMe string) ChannelOption {
	return func(c *channel) {
		c.managedByMe = managedByMe
	}
}

func WithChannelMaxResults(maxResults int64) ChannelOption {
	return func(c *channel) {
		c.maxResults = maxResults
	}
}

func WithChannelMine(mine string) ChannelOption {
	return func(c *channel) {
		c.mine = mine
	}
}

func WithChannelMySubscribers(mySubscribers string) ChannelOption {
	return func(c *channel) {
		c.mySubscribers = mySubscribers
	}
}

func WithChannelOnBehalfOfContentOwner(contentOwner string) ChannelOption {
	return func(c *channel) {
		c.onBehalfOfContentOwner = contentOwner
	}
}

func WithChannelCountry(country string) ChannelOption {
	return func(c *channel) {
		c.country = country
	}
}

func WithChannelCustomUrl(url string) ChannelOption {
	return func(c *channel) {
		c.customUrl = url
	}
}

func WithChannelDefaultLanguage(language string) ChannelOption {
	return func(c *channel) {
		c.defaultLanguage = language
	}
}

func WithChannelDescription(desc string) ChannelOption {
	return func(c *channel) {
		c.description = desc
	}
}

func WithChannelTitle(title string) ChannelOption {
	return func(c *channel) {
		c.title = title
	}
}
