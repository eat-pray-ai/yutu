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
	errGetMembershipsLevel = errors.New("failed to get memberships level")
)

type membershipsLevel struct{}

type MembershipsLevel interface {
	List([]string, string)
	get([]string) []*youtube.MembershipsLevel
}

type MembershipsLevelOption func(*membershipsLevel)

func NewMembershipsLevel(opts ...MembershipsLevelOption) MembershipsLevel {
	m := &membershipsLevel{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *membershipsLevel) get(parts []string) []*youtube.MembershipsLevel {
	call := service.MembershipsLevels.List(parts)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetMembershipsLevel, err))
	}

	return response.Items
}

func (m *membershipsLevel) List(parts []string, output string) {
	membershipsLevels := m.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(membershipsLevels)
	case "yaml":
		utils.PrintYAML(membershipsLevels)
	default:
		fmt.Println("id\tdisplayName")
		for _, membershipsLevel := range membershipsLevels {
			fmt.Printf(
				"%v\t%v\n", membershipsLevel.Id,
				membershipsLevel.Snippet.LevelDetails.DisplayName,
			)
		}
	}
}
