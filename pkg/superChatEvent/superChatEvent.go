package superChatEvent

import (
	"errors"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service              *youtube.Service
	errGetSuperChatEvent = errors.New("failed to get super chat event")
)

type superChatEvent struct {
	Hl         string `yaml:"hl" json:"hl"`
	MaxResults int64  `yaml:"max_results" json:"max_results"`
}

type SuperChatEvent interface {
	Get([]string) ([]*youtube.SuperChatEvent, error)
	List([]string, string, string, io.Writer) error
}

type Option func(*superChatEvent)

func NewSuperChatEvent(opts ...Option) SuperChatEvent {
	s := &superChatEvent{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *superChatEvent) Get(parts []string) ([]*youtube.SuperChatEvent, error) {
	call := service.SuperChatEvents.List(parts)
	if s.Hl != "" {
		call = call.Hl(s.Hl)
	}
	call = call.MaxResults(s.MaxResults)

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetSuperChatEvent, err)
	}

	return res.Items, nil
}

func (s *superChatEvent) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	events, err := s.Get(parts)
	if err != nil {
		return errors.Join(errGetSuperChatEvent, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(events, jpath, writer)
	case "yaml":
		utils.PrintYAML(events, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Amount", "Comment", "Supporter"})
		for _, event := range events {
			tb.AppendRow(
				table.Row{
					event.Id, event.Snippet.DisplayString, event.Snippet.CommentText,
					event.Snippet.SupporterDetails.DisplayName,
				},
			)
		}
	}
	return nil
}

func WithHl(hl string) Option {
	return func(s *superChatEvent) {
		s.Hl = hl
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *superChatEvent) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		s.MaxResults = maxResults
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *superChatEvent) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
