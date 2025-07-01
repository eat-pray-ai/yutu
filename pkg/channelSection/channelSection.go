package channelSection

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service                 *youtube.Service
	errGetChannelSection    = errors.New("failed to get channel section")
	errDeleteChannelSection = errors.New("failed to delete channel section")
)

type channelSection struct {
	IDs                    []string `yaml:"ids" json:"ids"`
	ChannelId              string   `yaml:"channel_id" json:"channel_id"`
	Hl                     string   `yaml:"hl" json:"hl"`
	Mine                   *bool    `yaml:"mine" json:"mine"`
	OnBehalfOfContentOwner string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
}

type ChannelSection interface {
	Get([]string) ([]*youtube.ChannelSection, error)
	List([]string, string, string, io.Writer) error
	Delete(writer io.Writer) error
	// Update()
	// Insert()
}

type Option func(*channelSection)

func NewChannelSection(opts ...Option) ChannelSection {
	cs := &channelSection{}

	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

func (cs *channelSection) Get(parts []string) (
	[]*youtube.ChannelSection, error,
) {
	call := service.ChannelSections.List(parts)
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

func (cs *channelSection) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	channelSections, err := cs.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(channelSections, jpath, writer)
	case "yaml":
		utils.PrintYAML(channelSections, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Channel ID", "Title"})
		for _, chs := range channelSections {
			tb.AppendRow(table.Row{chs.Id, chs.Snippet.ChannelId, chs.Snippet.Title})
		}
	}
	return nil
}

func (cs *channelSection) Delete(writer io.Writer) error {
	for _, id := range cs.IDs {
		call := service.ChannelSections.Delete(id)
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

func WithIDs(ids []string) Option {
	return func(cs *channelSection) {
		cs.IDs = ids
	}
}

func WithChannelId(channelId string) Option {
	return func(cs *channelSection) {
		cs.ChannelId = channelId
	}
}

func WithHl(hl string) Option {
	return func(cs *channelSection) {
		cs.Hl = hl
	}
}

func WithMine(mine *bool) Option {
	return func(cs *channelSection) {
		if mine != nil {
			cs.Mine = mine
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cs *channelSection) {
		cs.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *channelSection) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
