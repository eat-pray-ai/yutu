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
	CategoryId             string `yaml:"category_id" json:"category_id"`
	ForHandle              string `yaml:"for_handle" json:"for_handle"`
	ForUsername            string `yaml:"for_username" json:"for_username"`
	Hl                     string `yaml:"hl" json:"hl"`
	ID                     string `yaml:"id" json:"id"`
	ManagedByMe            *bool  `yaml:"managed_by_me" json:"managed_by_me"`
	MaxResults             int64  `yaml:"max_results" json:"max_results"`
	Mine                   *bool  `yaml:"mine" json:"mine"`
	MySubscribers          *bool  `yaml:"my_subscribers" json:"my_subscribers"`
	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`

	Country         string `yaml:"country" json:"country"`
	CustomUrl       string `yaml:"custom_url" json:"custom_url"`
	DefaultLanguage string `yaml:"default_language" json:"default_language"`
	Description     string `yaml:"description" json:"description"`
	Title           string `yaml:"title" json:"title"`
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
	if c.CategoryId != "" {
		call = call.CategoryId(c.CategoryId)
	}

	if c.ForHandle != "" {
		call = call.ForUsername(c.ForHandle)
	}

	if c.ForUsername != "" {
		call = call.ForUsername(c.ForUsername)
	}

	if c.Hl != "" {
		call = call.Hl(c.Hl)
	}

	if c.ID != "" {
		call = call.Id(c.ID)
	}

	if c.ManagedByMe != nil {
		call = call.ManagedByMe(*c.ManagedByMe)
	}

	call = call.MaxResults(c.MaxResults)

	if c.Mine != nil {
		call = call.Mine(*c.Mine)
	}

	if c.MySubscribers != nil {
		call = call.MySubscribers(*c.MySubscribers)
	}

	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errGetChannel, err), c.ID)
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
	if c.Title != "" {
		channel.Snippet.Title = c.Title
	}
	if c.Description != "" {
		channel.Snippet.Description = c.Description
	}

	call := service.Channels.Update(parts, channel)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errUpdateChannel, err), c.ID)
	}
	utils.PrintYAML(res)
}

func WithCategoryId(categoryId string) Option {
	return func(c *channel) {
		c.CategoryId = categoryId
	}
}

func WithForHandle(handle string) Option {
	return func(c *channel) {
		c.ForHandle = handle
	}
}

func WithForUsername(username string) Option {
	return func(c *channel) {
		c.ForUsername = username
	}
}

func WithHl(hl string) Option {
	return func(c *channel) {
		c.Hl = hl
	}
}

func WithID(id string) Option {
	return func(c *channel) {
		c.ID = id
	}
}

func WithChannelManagedByMe(managedByMe bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.ManagedByMe = &managedByMe
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *channel) {
		c.MaxResults = maxResults
	}
}

func WithMine(mine bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.Mine = &mine
		}
	}
}

func WithMySubscribers(mySubscribers bool, changed bool) Option {
	return func(c *channel) {
		if changed {
			c.MySubscribers = &mySubscribers
		}
	}
}

func WithOnBehalfOfContentOwner(contentOwner string) Option {
	return func(c *channel) {
		c.OnBehalfOfContentOwner = contentOwner
	}
}

func WithCountry(country string) Option {
	return func(c *channel) {
		c.Country = country
	}
}

func WithCustomUrl(url string) Option {
	return func(c *channel) {
		c.CustomUrl = url
	}
}

func WithDefaultLanguage(language string) Option {
	return func(c *channel) {
		c.DefaultLanguage = language
	}
}

func WithDescription(desc string) Option {
	return func(c *channel) {
		c.Description = desc
	}
}

func WithTitle(title string) Option {
	return func(c *channel) {
		c.Title = title
	}
}

func WithService(svc *youtube.Service) Option {
	return func(c *channel) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
