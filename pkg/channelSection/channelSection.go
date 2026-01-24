// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetChannelSection    = errors.New("failed to get channel section")
	errDeleteChannelSection = errors.New("failed to delete channel section")
)

type ChannelSection struct {
	service                *youtube.Service
	IDs                    []string `yaml:"ids" json:"ids"`
	ChannelId              string   `yaml:"channel_id" json:"channel_id"`
	Hl                     string   `yaml:"hl" json:"hl"`
	Mine                   *bool    `yaml:"mine" json:"mine"`
	OnBehalfOfContentOwner string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`

	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
}

type IChannelSection[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Delete(io.Writer) error
	preRun()
	// Update()
	// Insert()
}

type Option func(*ChannelSection)

func NewChannelSection(opts ...Option) IChannelSection[youtube.ChannelSection] {
	cs := &ChannelSection{}

	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

func (cs *ChannelSection) Get() (
	[]*youtube.ChannelSection, error,
) {
	cs.preRun()
	call := cs.service.ChannelSections.List(cs.Parts)
	if len(cs.IDs) > 0 {
		call = call.Id(cs.IDs...)
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
	cs.preRun()
	for _, id := range cs.IDs {
		call := cs.service.ChannelSections.Delete(id)
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

func (cs *ChannelSection) preRun() {
	if cs.service == nil {
		cs.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func WithIDs(ids []string) Option {
	return func(cs *ChannelSection) {
		cs.IDs = ids
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

func WithParts(parts []string) Option {
	return func(cs *ChannelSection) {
		cs.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(cs *ChannelSection) {
		cs.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(cs *ChannelSection) {
		cs.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(cs *ChannelSection) {
		cs.service = svc
	}
}
