// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannelSection    = errors.New("failed to get channel section")
	errDeleteChannelSection = errors.New("failed to delete channel section")
)

type ChannelSection struct {
	*common.Fields
	Ids                    []string `yaml:"ids" json:"ids"`
	ChannelId              string   `yaml:"channel_id" json:"channel_id"`
	Hl                     string   `yaml:"hl" json:"hl"`
	Mine                   *bool    `yaml:"mine" json:"mine"`
	OnBehalfOfContentOwner string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
}

type IChannelSection[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Delete(io.Writer) error
	// Update()
	// Insert()
}

type Option func(*ChannelSection)

func NewChannelSection(opts ...Option) IChannelSection[youtube.ChannelSection] {
	cs := &ChannelSection{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

func (cs *ChannelSection) GetFields() *common.Fields {
	return cs.Fields
}

func (cs *ChannelSection) Get() (
	[]*youtube.ChannelSection, error,
) {
	cs.EnsureService()
	call := cs.Service.ChannelSections.List(cs.Parts)
	if len(cs.Ids) > 0 {
		call = call.Id(cs.Ids...)
	}
	if cs.ChannelId != "" {
		call = call.ChannelId(cs.ChannelId)
	}
	if cs.Hl != "" {
		call = call.Hl(cs.Hl)
	}
	if cs.Mine != nil {
		call = call.Mine(*cs.Mine)
	}
	if cs.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cs.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetChannelSection, err)
	}
	return res.Items, nil
}

func (cs *ChannelSection) List(writer io.Writer) error {
	channelSections, err := cs.Get()
	if err != nil {
		return err
	}

	switch cs.Output {
	case "json":
		utils.PrintJSON(channelSections, cs.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(channelSections, cs.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Channel ID", "Title"})
		for _, chs := range channelSections {
			tb.AppendRow(table.Row{chs.Id, chs.Snippet.ChannelId, chs.Snippet.Title})
		}
	}
	return nil
}

func (cs *ChannelSection) Delete(writer io.Writer) error {
	cs.EnsureService()
	for _, id := range cs.Ids {
		call := cs.Service.ChannelSections.Delete(id)
		if cs.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(cs.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteChannelSection, err)
		}

		_, _ = fmt.Fprintf(writer, "Channel section %s deleted\n", id)
	}
	return nil
}

func WithIds(ids []string) Option {
	return func(cs *ChannelSection) {
		cs.Ids = ids
	}
}

func WithChannelId(channelId string) Option {
	return func(cs *ChannelSection) {
		cs.ChannelId = channelId
	}
}

func WithHl(hl string) Option {
	return func(cs *ChannelSection) {
		cs.Hl = hl
	}
}

func WithMine(mine *bool) Option {
	return func(cs *ChannelSection) {
		if mine != nil {
			cs.Mine = mine
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cs *ChannelSection) {
		cs.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

var (
	WithParts    = common.WithParts[*ChannelSection]
	WithOutput   = common.WithOutput[*ChannelSection]
	WithJsonpath = common.WithJsonpath[*ChannelSection]
	WithService  = common.WithService[*ChannelSection]
)
