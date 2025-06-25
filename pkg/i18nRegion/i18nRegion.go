package i18nRegion

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
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
		_, _ = fmt.Fprintln(writer, "ID\tGl\tName")
		for _, i18nRegion := range i18nRegions {
			_, _ = fmt.Fprintf(
				writer, "%s\t%s\t%s\n",
				i18nRegion.Id, i18nRegion.Snippet.Gl, i18nRegion.Snippet.Name,
			)
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
