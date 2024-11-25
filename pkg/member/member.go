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
	MemberChannelId  string `yaml:"member_channel_id" json:"member_channel_id"`
	HasAccessToLevel string `yaml:"has_access_to_level" json:"has_access_to_level"`
	MaxResults       int64  `yaml:"max_results" json:"max_results"`
	Mode             string `yaml:"mode" json:"mode"`
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
		utils.PrintJSON(m)
		log.Fatalln(errors.Join(errGetMember, err))
	}

	return res.Items
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
