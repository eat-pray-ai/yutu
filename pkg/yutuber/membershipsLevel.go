package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	errGetMembershipsLevel error = errors.New("failed to get memberships level")
)

type MembershipsLevel struct{}

type MembershipsLevelService interface {
	List([]string, string)
	get([]string) []*youtube.MembershipsLevel
}

type MembershipsLevelOption func(*MembershipsLevel)

func NewMembershipsLevel(opts ...MembershipsLevelOption) *MembershipsLevel {
	m := &MembershipsLevel{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *MembershipsLevel) get(parts []string) []*youtube.MembershipsLevel {
	call := service.MembershipsLevels.List(parts)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetMembershipsLevel, err))
	}

	return response.Items
}

func (m *MembershipsLevel) List(parts []string, output string) {
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
