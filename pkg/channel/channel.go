// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannel    = errors.New("failed to get channel")
	errUpdateChannel = errors.New("failed to update channel")
)

type Channel struct {
	*pkg.DefaultFields
	CategoryId             string   `yaml:"category_id" json:"category_id"`
	ForHandle              string   `yaml:"for_handle" json:"for_handle"`
	ForUsername            string   `yaml:"for_username" json:"for_username"`
	Hl                     string   `yaml:"hl" json:"hl"`
	Ids                    []string `yaml:"ids" json:"ids"`
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

type IChannel[T youtube.Channel] interface {
	List(io.Writer) error
	Update(io.Writer) error
	Get() ([]*T, error)
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(*Channel)

func NewChannel(opts ...Option) IChannel[youtube.Channel] {
	c := &Channel{DefaultFields: &pkg.DefaultFields{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Channel) GetDefaultFields() *pkg.DefaultFields {
	return c.DefaultFields
}

func (c *Channel) Get() ([]*youtube.Channel, error) {
	c.preRun()
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

	var items []*youtube.Channel
	pageToken := ""
	for c.MaxResults > 0 {
		call = call.MaxResults(min(c.MaxResults, pkg.PerPage))
		c.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetChannel, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (c *Channel) List(writer io.Writer) error {
	channels, err := c.Get()
	if err != nil && channels == nil {
		return err
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(channels, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(channels, c.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Title", "Country"})
		for _, channel := range channels {
			tb.AppendRow(
				table.Row{channel.Id, channel.Snippet.Title, channel.Snippet.Country},
			)
		}
	}
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

	switch c.Output {
	case "json":
		utils.PrintJSON(res, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Channel updated: %s\n", res.Id)
	}
	return nil
}

func (c *Channel) preRun() {
	if c.Service == nil {
		c.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
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

func WithHl(hl string) Option {
	return func(c *Channel) {
		c.Hl = hl
	}
}

func WithIds(ids []string) Option {
	return func(c *Channel) {
		c.Ids = ids
	}
}

func WithChannelManagedByMe(managedByMe *bool) Option {
	return func(c *Channel) {
		if managedByMe != nil {
			c.ManagedByMe = managedByMe
		}
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *Channel) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		c.MaxResults = maxResults
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

func WithOnBehalfOfContentOwner(contentOwner string) Option {
	return func(c *Channel) {
		c.OnBehalfOfContentOwner = contentOwner
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
	WithParts    = pkg.WithParts[*Channel]
	WithOutput   = pkg.WithOutput[*Channel]
	WithJsonpath = pkg.WithJsonpath[*Channel]
	WithService  = pkg.WithService[*Channel]
)
