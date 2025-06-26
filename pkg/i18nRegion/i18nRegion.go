package i18nRegion

import (
	"errors"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service          *youtube.Service
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type i18nRegion struct {
	Hl string `yaml:"hl" json:"hl"`
}

type I18nRegion interface {
	Get([]string) ([]*youtube.I18nRegion, error)
	List([]string, string, io.Writer) error
}

type Option func(*i18nRegion)

func NewI18nRegion(opts ...Option) I18nRegion {
	i := &i18nRegion{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *i18nRegion) Get(parts []string) ([]*youtube.I18nRegion, error) {
	call := service.I18nRegions.List(parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetI18nRegion, err)
	}

	return res.Items, nil
}

func (i *i18nRegion) List(
	parts []string, output string, writer io.Writer,
) error {
	i18nRegions, err := i.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(i18nRegions, writer)
	case "yaml":
		utils.PrintYAML(i18nRegions, writer)
	default:
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Gl", "Name"})
		for _, region := range i18nRegions {
			tb.AppendRow(table.Row{region.Id, region.Snippet.Gl, region.Snippet.Name})
		}
	}
	return nil
}

func WithHl(hl string) Option {
	return func(i *i18nRegion) {
		i.Hl = hl
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *i18nRegion) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
