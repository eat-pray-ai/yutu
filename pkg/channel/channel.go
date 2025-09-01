package channel

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"

	"google.golang.org/api/youtube/v3"
)

var (
	service          *youtube.Service
	errGetChannel    = errors.New("failed to get channel")
	errUpdateChannel = errors.New("failed to update channel")
)

type channel struct {
	CategoryId             string   `yaml:"category_id" json:"category_id"`
	ForHandle              string   `yaml:"for_handle" json:"for_handle"`
	ForUsername            string   `yaml:"for_username" json:"for_username"`
	Hl                     string   `yaml:"hl" json:"hl"`
	IDs                    []string `yaml:"ids" json:"ids"`
	ManagedByMe            *bool    `yaml:"managed_by_me" json:"managed_by_me"`
	MaxResults             int64    `yaml:"max_results" json:"max_results"`
	Mine                   *bool    `yaml:"mine" json:"mine"`
	MySubscribers          *bool    `yaml:"my_subscribers" json:"my_subscribers"`
	OnBehalfOfContentOwner string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`

	Country         string `yaml:"country" json:"country"`
	CustomUrl       string `yaml:"custom_url" json:"custom_url"`
	DefaultLanguage string `yaml:"default_language" json:"default_language"`
	Description     string `yaml:"description" json:"description"`
	Title           string `yaml:"title" json:"title"`
}

type Channel interface {
	List([]string, string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Get([]string) ([]*youtube.Channel, error)
}

type Option func(*channel)

func NewChannel(opts ...Option) Channel {
	c := &channel{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *channel) Get(parts []string) ([]*youtube.Channel, error) {
	call := service.Channels.List(parts)
	if c.CategoryId != "" {
		call = call.CategoryId(c.CategoryId)
	}

	if c.ForHandle != "" {
		call = call.ForHandle(c.ForHandle)
	}

	if c.ForUsername != "" {
		call = call.ForUsername(c.ForUsername)
	}

	if c.Hl != "" {
		call = call.Hl(c.Hl)
	}

	if len(c.IDs) > 0 {
		call = call.Id(c.IDs...)
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
		return nil, errors.Join(errGetChannel, err)
	}

	return res.Items, nil
}

func (c *channel) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	channels, err := c.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(channels, jpath, writer)
	case "yaml":
		utils.PrintYAML(channels, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Title", "Country"})
		for _, channel := range channels {
			tb.AppendRow(
				table.Row{channel.Id, channel.Snippet.Title, channel.Snippet.Country},
			)
		}
	}
	return nil
}

func (c *channel) Update(output string, jpath string, writer io.Writer) error {
	parts := []string{"snippet"}
	channels, err := c.Get(parts)
	if err != nil {
		return errors.Join(errUpdateChannel, err)
	}
	if len(channels) == 0 {
		return errGetChannel
	}

	cha := channels[0]
	if c.Country != "" {
		cha.Snippet.Country = c.Country
	}
	if c.CustomUrl != "" {
		cha.Snippet.CustomUrl = c.CustomUrl
	}
	if c.DefaultLanguage != "" {
		cha.Snippet.DefaultLanguage = c.DefaultLanguage
	}
	if c.Description != "" {
		cha.Snippet.Description = c.Description
	}
	if c.Title != "" {
		cha.Snippet.Title = c.Title
	}

	call := service.Channels.Update(parts, cha)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateChannel, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Channel updated: %s\n", res.Id)
	}
	return nil
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

func WithIDs(ids []string) Option {
	return func(c *channel) {
		c.IDs = ids
	}
}

func WithChannelManagedByMe(managedByMe *bool) Option {
	return func(c *channel) {
		if managedByMe != nil {
			c.ManagedByMe = managedByMe
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *channel) {
		if maxResults <= 0 {
			maxResults = 1
		}
		c.MaxResults = maxResults
	}
}

func WithMine(mine *bool) Option {
	return func(c *channel) {
		if mine != nil {
			c.Mine = mine
		}
	}
}

func WithMySubscribers(mySubscribers *bool) Option {
	return func(c *channel) {
		if mySubscribers != nil {
			c.MySubscribers = mySubscribers
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
	return func(_ *channel) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			).GetService()
		}
		service = svc
	}
}
