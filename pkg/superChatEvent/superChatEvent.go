package superChatEvent

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"io"
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
	List([]string, string, io.Writer) error
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
	if s.MaxResults <= 0 {
		s.MaxResults = 1
	}
	call = call.MaxResults(s.MaxResults)

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetSuperChatEvent, err)
	}

	return res.Items, nil
}

func (s *superChatEvent) List(
	parts []string, output string, writer io.Writer,
) error {
	events, err := s.Get(parts)
	if err != nil {
		return errors.Join(errGetSuperChatEvent, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(events, writer)
	case "yaml":
		utils.PrintYAML(events, writer)
	default:
		_, _ = fmt.Fprintln(writer, "ID\tAmount")
		for _, event := range events {
			_, _ = fmt.Fprintf(
				writer, "%v\t%v\n", event.Id, event.Snippet.AmountMicros,
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
		s.MaxResults = maxResults
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *superChatEvent) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
