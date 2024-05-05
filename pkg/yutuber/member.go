package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	errGetMember = errors.New("failed to get member")
)

type member struct {
	memberChannelId string
}

type Member interface {
	List([]string, string)
	get([]string) []*youtube.Member
}

type MemberOption func(*member)

func NewMember(opts ...MemberOption) Member {
	m := &member{}
	service = auth.NewY2BService()

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

func WithMemberChannelId(channelId string) MemberOption {
	return func(m *member) {
		m.memberChannelId = channelId
	}
}
