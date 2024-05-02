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
	errGetI18nLanguage = errors.New("failed to get i18n language")
)

type i18nLanguage struct{}

type I18nLanguage interface {
	get(parts []string) []*youtube.I18nLanguage
	List(parts []string, output string)
}

type i18nLanguageOption func(*i18nLanguage)

func NewI18nLanguage(opts ...i18nLanguageOption) I18nLanguage {
	service = auth.NewY2BService()
	i := &i18nLanguage{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *i18nLanguage) get(parts []string) []*youtube.I18nLanguage {
	call := service.I18nLanguages.List(parts)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetI18nLanguage, err))
	}

	return response.Items
}

func (i *i18nLanguage) List(parts []string, output string) {
	i18nLanguages := i.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(i18nLanguages)
	case "yaml":
		utils.PrintYAML(i18nLanguages)
	default:
		fmt.Println("id\thl\tname")
		for _, i18nLanguage := range i18nLanguages {
			fmt.Printf(
				"%v\t%v\t%v\n",
				i18nLanguage.Id, i18nLanguage.Snippet.Hl, i18nLanguage.Snippet.Name,
			)
		}
	}
}
