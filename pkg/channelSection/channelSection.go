// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannelSection    = errors.New("failed to get channel section")
	errDeleteChannelSection = errors.New("failed to delete channel section")
)

type ChannelSection struct {
	*common.Fields
	Mine *bool `yaml:"mine" json:"mine,omitempty"`
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
	if err := cs.EnsureService(); err != nil {
		return nil, err
	}
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

	common.PrintList(cs.Output, channelSections, writer, table.Row{"ID", "Channel ID", "Title"}, func(s *youtube.ChannelSection) table.Row {
		return table.Row{s.Id, s.Snippet.ChannelId, s.Snippet.Title}
	})
	return nil
}

func (cs *ChannelSection) Delete(writer io.Writer) error {
	if err := cs.EnsureService(); err != nil {
		return err
	}
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

func WithMine(mine *bool) Option {
	return func(cs *ChannelSection) {
		if mine != nil {
			cs.Mine = mine
		}
	}
}

var (
	WithParts     = common.WithParts[*ChannelSection]
	WithOutput    = common.WithOutput[*ChannelSection]
	WithService   = common.WithService[*ChannelSection]
	WithIds       = common.WithIds[*ChannelSection]
	WithHl        = common.WithHl[*ChannelSection]
	WithChannelId = common.WithChannelId[*ChannelSection]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*ChannelSection]
)
