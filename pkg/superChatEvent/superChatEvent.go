package superChatEvent

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
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
	get([]string) []*youtube.SuperChatEvent
	List([]string, string)
}

type Option func(*superChatEvent)

func NewSuperChatEvent(opts ...Option) SuperChatEvent {
	s := &superChatEvent{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *superChatEvent) get(parts []string) []*youtube.SuperChatEvent {
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
		utils.PrintJSON(s)
		log.Fatalln(errors.Join(errGetSuperChatEvent, err))
	}

	return res.Items
}

func (s *superChatEvent) List(parts []string, output string) {
	events := s.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(events)
	case "yaml":
		utils.PrintYAML(events)
	default:
		fmt.Println("ID\tAmount")
		for _, event := range events {
			fmt.Printf("%v\t%v\n", event.Id, event.Snippet.AmountMicros)
		}
	}
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
