package member

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service      *youtube.Service
	errGetMember = errors.New("failed to get member")
)

type member struct {
	memberChannelId  string
	hasAccessToLevel string
	maxResults       int64
	mode             string
}

type Member interface {
	List([]string, string)
	get([]string) []*youtube.Member
}

type Option func(*member)

func NewMember(opts ...Option) Member {
	m := &member{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *member) get(parts []string) []*youtube.Member {
	call := service.Members.List(parts)
	if m.memberChannelId != "" {
		call = call.FilterByMemberChannelId(m.memberChannelId)
	}
	if m.hasAccessToLevel != "" {
		call = call.HasAccessToLevel(m.hasAccessToLevel)
	}
	call = call.MaxResults(m.maxResults)
	if m.mode != "" {
		call = call.Mode(m.mode)
	}

	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetMember, err))
	}

	return response.Items
}

func (m *member) List(parts []string, output string) {
	members := m.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(members)
	case "yaml":
		utils.PrintYAML(members)
	default:
		fmt.Println("channelId\tdisplayName")
		for _, member := range members {
			fmt.Printf(
				"%s\t%s\n",
				member.Snippet.MemberDetails.ChannelId,
				member.Snippet.MemberDetails.DisplayName,
			)
		}
	}
}

func WithMemberChannelId(channelId string) Option {
	return func(m *member) {
		m.memberChannelId = channelId
	}
}

func WithHasAccessToLevel(level string) Option {
	return func(m *member) {
		m.hasAccessToLevel = level
	}
}

func WithMaxResults(results int64) Option {
	return func(m *member) {
		m.maxResults = results
	}
}

func WithMode(mode string) Option {
	return func(m *member) {
		m.mode = mode
	}
}

func WithService() Option {
	return func(m *member) {
		service = auth.NewY2BService()
	}
}
