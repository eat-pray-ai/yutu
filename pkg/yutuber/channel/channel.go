package channel

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
	managedByMe            *bool
	maxResults             int64
	mine                   *bool
	mySubscribers          *bool
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

type Option func(*channel)

func NewChannel(opts ...Option) Channel {
	c := &channel{}

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

	if c.managedByMe != nil {
		call = call.ManagedByMe(*c.managedByMe)
	}

	call = call.MaxResults(c.maxResults)

	if c.mine != nil {
		call = call.Mine(*c.mine)
	}

	if c.mySubscribers != nil {
		call = call.MySubscribers(*c.mySubscribers)
	}

	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetChannel, err), c.id)
	}

	return res.Items
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
	utils.PrintYAML(res)
}

func WithCategoryId(categoryId string) Option {
	return func(c *channel) {
		c.categoryId = categoryId
	}
}

func WithForHandle(handle string) Option {
	return func(c *channel) {
		c.forHandle = handle
	}
}

func WithForUsername(username string) Option {
	return func(c *channel) {
		c.forUsername = username
	}
}

func WithHl(hl string) Option {
	return func(c *channel) {
		c.hl = hl
	}
}

func WithId(id string) Option {
	return func(c *channel) {
		c.id = id
	}
}

func WithChannelManagedByMe(managedByMe bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.managedByMe = &managedByMe
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *channel) {
		c.maxResults = maxResults
	}
}

func WithMine(mine bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.mine = &mine
		}
	}
}

func WithMySubscribers(mySubscribers bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.mySubscribers = &mySubscribers
		}
	}
}

func WithOnBehalfOfContentOwner(contentOwner string) Option {
	return func(c *channel) {
		c.onBehalfOfContentOwner = contentOwner
	}
}

func WithCountry(country string) Option {
	return func(c *channel) {
		c.country = country
	}
}

func WithCustomUrl(url string) Option {
	return func(c *channel) {
		c.customUrl = url
	}
}

func WithDefaultLanguage(language string) Option {
	return func(c *channel) {
		c.defaultLanguage = language
	}
}

func WithDescription(desc string) Option {
	return func(c *channel) {
		c.description = desc
	}
}

func WithTitle(title string) Option {
	return func(c *channel) {
		c.title = title
	}
}

func WithService() Option {
	return func(c *channel) {
		service = auth.NewY2BService()
	}
}
