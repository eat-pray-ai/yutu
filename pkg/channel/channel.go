// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannel    = errors.New("failed to get channel")
	errUpdateChannel = errors.New("failed to update channel")
)

type Channel struct {
	*common.Fields
	CategoryId      string `yaml:"category_id" json:"category_id,omitempty"`
	ForHandle       string `yaml:"for_handle" json:"for_handle,omitempty"`
	ForUsername     string `yaml:"for_username" json:"for_username,omitempty"`
	ManagedByMe     *bool  `yaml:"managed_by_me" json:"managed_by_me,omitempty"`
	Mine            *bool  `yaml:"mine" json:"mine,omitempty"`
	MySubscribers   *bool  `yaml:"my_subscribers" json:"my_subscribers,omitempty"`
	Country         string `yaml:"country" json:"country,omitempty"`
	CustomUrl       string `yaml:"custom_url" json:"custom_url,omitempty"`
	DefaultLanguage string `yaml:"default_language" json:"default_language,omitempty"`
	Description     string `yaml:"description" json:"description,omitempty"`
	Title           string `yaml:"title" json:"title,omitempty"`
}

type IChannel[T youtube.Channel] interface {
	List(io.Writer) error
	Update(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*Channel)

func NewChannel(opts ...Option) IChannel[youtube.Channel] {
	c := &Channel{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Channel) Get() ([]*youtube.Channel, error) {
	if err := c.EnsureService(); err != nil {
		return nil, err
	}
	call := c.Service.Channels.List(c.Parts)
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
	if len(c.Ids) > 0 {
		call = call.Id(c.Ids...)
	}
	if c.ManagedByMe != nil {
		call = call.ManagedByMe(*c.ManagedByMe)
	}
	if c.Mine != nil {
		call = call.Mine(*c.Mine)
	}
	if c.MySubscribers != nil {
		call = call.MySubscribers(*c.MySubscribers)
	}
	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	return common.Paginate(c.Fields, call, func(r *youtube.ChannelListResponse) ([]*youtube.Channel, string) {
		return r.Items, r.NextPageToken
	}, errGetChannel)
}

func (c *Channel) List(writer io.Writer) error {
	channels, err := c.Get()
	if err != nil && channels == nil {
		return err
	}

	common.PrintList(c.Output, channels, writer, table.Row{"ID", "Title", "Country"}, func(ch *youtube.Channel) table.Row {
		title := ""
		country := ""
		if ch.Snippet != nil {
			title = ch.Snippet.Title
			country = ch.Snippet.Country
		}
		return table.Row{ch.Id, title, country}
	})
	return err
}

func (c *Channel) Update(writer io.Writer) error {
	c.Parts = []string{"snippet"}
	channels, err := c.Get()
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

	call := c.Service.Channels.Update(c.Parts, cha)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateChannel, err)
	}

	common.PrintResult(c.Output, res, writer, "Channel updated: %s\n", res.Id)
	return nil
}

func WithCategoryId(categoryId string) Option {
	return func(c *Channel) {
		c.CategoryId = categoryId
	}
}

func WithForHandle(handle string) Option {
	return func(c *Channel) {
		c.ForHandle = handle
	}
}

func WithForUsername(username string) Option {
	return func(c *Channel) {
		c.ForUsername = username
	}
}

func WithChannelManagedByMe(managedByMe *bool) Option {
	return func(c *Channel) {
		if managedByMe != nil {
			c.ManagedByMe = managedByMe
		}
	}
}

func WithMine(mine *bool) Option {
	return func(c *Channel) {
		if mine != nil {
			c.Mine = mine
		}
	}
}

func WithMySubscribers(mySubscribers *bool) Option {
	return func(c *Channel) {
		if mySubscribers != nil {
			c.MySubscribers = mySubscribers
		}
	}
}

func WithCountry(country string) Option {
	return func(c *Channel) {
		c.Country = country
	}
}

func WithCustomUrl(url string) Option {
	return func(c *Channel) {
		c.CustomUrl = url
	}
}

func WithDefaultLanguage(language string) Option {
	return func(c *Channel) {
		c.DefaultLanguage = language
	}
}

func WithDescription(desc string) Option {
	return func(c *Channel) {
		c.Description = desc
	}
}

func WithTitle(title string) Option {
	return func(c *Channel) {
		c.Title = title
	}
}

var (
	WithParts      = common.WithParts[*Channel]
	WithOutput     = common.WithOutput[*Channel]
	WithService    = common.WithService[*Channel]
	WithIds        = common.WithIds[*Channel]
	WithMaxResults = common.WithMaxResults[*Channel]
	WithHl         = common.WithHl[*Channel]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*Channel]
)
