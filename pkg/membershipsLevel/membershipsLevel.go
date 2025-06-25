package membershipsLevel

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service                *youtube.Service
	errGetMembershipsLevel = errors.New("failed to get memberships level")
)

type membershipsLevel struct{}

type MembershipsLevel interface {
	List([]string, string, io.Writer) error
	Get([]string) ([]*youtube.MembershipsLevel, error)
}

type Option func(*membershipsLevel)

func NewMembershipsLevel(opts ...Option) MembershipsLevel {
	m := &membershipsLevel{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *membershipsLevel) Get(parts []string) (
	[]*youtube.MembershipsLevel, error,
) {
	call := service.MembershipsLevels.List(parts)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetMembershipsLevel, err)
	}

	return res.Items, nil
}

func (m *membershipsLevel) List(
	parts []string, output string, writer io.Writer,
) error {
	membershipsLevels, err := m.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(membershipsLevels, writer)
	case "yaml":
		utils.PrintYAML(membershipsLevels, writer)
	default:
		_, _ = fmt.Fprintln(writer, "ID\tDisplayName")
		for _, membershipsLevel := range membershipsLevels {
			_, _ = fmt.Fprintf(
				writer, "%v\t%v\n", membershipsLevel.Id,
				membershipsLevel.Snippet.LevelDetails.DisplayName,
			)
		}
	}
	return nil
}

func WithService(svc *youtube.Service) Option {
	return func(_ *membershipsLevel) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
