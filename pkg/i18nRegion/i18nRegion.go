package i18nRegion

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service          *youtube.Service
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type i18nRegion struct {
	Hl string `yaml:"hl" json:"hl"`
}

type I18nRegion interface {
	get(parts []string) []*youtube.I18nRegion
	List(parts []string, output string)
}

type Option func(*i18nRegion)

func NewI18nRegion(opts ...Option) I18nRegion {
	i := &i18nRegion{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *i18nRegion) get(parts []string) []*youtube.I18nRegion {
	call := service.I18nRegions.List(parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(i)
		log.Fatalln(errors.Join(errGetI18nRegion, err))
	}

	return res.Items
}

func (i *i18nRegion) List(parts []string, output string) {
	i18nRegions := i.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(i18nRegions)
	case "yaml":
		utils.PrintYAML(i18nRegions)
	default:
		fmt.Println("ID\tgl\tname")
		for _, i18nRegion := range i18nRegions {
			fmt.Printf(
				"%v\t%v\t%v\n",
				i18nRegion.Id, i18nRegion.Snippet.Gl, i18nRegion.Snippet.Name,
			)
		}
	}
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
