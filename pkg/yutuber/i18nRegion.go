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
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type i18nRegion struct{}

type I18nRegion interface {
	get(parts []string) []*youtube.I18nRegion
	List(parts []string, output string)
}

type i18nRegionOption func(*i18nRegion)

func NewI18nRegion(opts ...i18nRegionOption) I18nRegion {
	service = auth.NewY2BService()
	i := &i18nRegion{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *i18nRegion) get(parts []string) []*youtube.I18nRegion {
	call := service.I18nRegions.List(parts)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetI18nRegion, err))
	}

	return response.Items
}

func (i *i18nRegion) List(parts []string, output string) {
	i18nRegions := i.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(i18nRegions)
	case "yaml":
		utils.PrintYAML(i18nRegions)
	default:
		fmt.Println("Id\tgl\tname")
		for _, i18nRegion := range i18nRegions {
			fmt.Printf(
				"%v\t%v\t%v\n",
				i18nRegion.Id, i18nRegion.Snippet.Gl, i18nRegion.Snippet.Name,
			)
		}
	}
}
