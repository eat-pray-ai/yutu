package member

import (
	"errors"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service      *youtube.Service
	errGetMember = errors.New("failed to get member")
)

type member struct {
	MemberChannelId  string `yaml:"member_channel_id" json:"member_channel_id"`
	HasAccessToLevel string `yaml:"has_access_to_level" json:"has_access_to_level"`
	MaxResults       int64  `yaml:"max_results" json:"max_results"`
	Mode             string `yaml:"mode" json:"mode"`
}

type Member interface {
	List([]string, string, string, io.Writer) error
	Get([]string) ([]*youtube.Member, error)
}

type Option func(*member)

func NewMember(opts ...Option) Member {
	m := &member{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *member) Get(parts []string) ([]*youtube.Member, error) {
	call := service.Members.List(parts)
	if m.MemberChannelId != "" {
		call = call.FilterByMemberChannelId(m.MemberChannelId)
	}
	if m.HasAccessToLevel != "" {
		call = call.HasAccessToLevel(m.HasAccessToLevel)
	}
	if m.MaxResults <= 0 {
		m.MaxResults = 1
	}
	call = call.MaxResults(m.MaxResults)
	if m.Mode != "" {
		call = call.Mode(m.Mode)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetMember, err)
	}

	return res.Items, nil
}

func (m *member) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	members, err := m.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(members, jpath, writer)
	case "yaml":
		utils.PrintYAML(members, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"Channel ID", "Display Name"})
		for _, member := range members {
			tb.AppendRow(
				table.Row{
					member.Snippet.MemberDetails.ChannelId,
					member.Snippet.MemberDetails.DisplayName,
				},
			)
		}
	}
	return nil
}

func WithMemberChannelId(channelId string) Option {
	return func(m *member) {
		m.MemberChannelId = channelId
	}
}

func WithHasAccessToLevel(level string) Option {
	return func(m *member) {
		m.HasAccessToLevel = level
	}
}

func WithMaxResults(results int64) Option {
	return func(m *member) {
		m.MaxResults = results
	}
}

func WithMode(mode string) Option {
	return func(m *member) {
		m.Mode = mode
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *member) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
