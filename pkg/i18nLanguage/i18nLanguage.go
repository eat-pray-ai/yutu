// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service            *youtube.Service
	errGetI18nLanguage = errors.New("failed to get i18n language")
)

type i18nLanguage struct {
	Hl string `yaml:"hl" json:"hl"`
}

type I18nLanguage[T youtube.I18nLanguage] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
}

type Option func(*i18nLanguage)

func NewI18nLanguage(opts ...Option) I18nLanguage[youtube.I18nLanguage] {
	i := &i18nLanguage{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *i18nLanguage) Get(parts []string) (
	[]*youtube.I18nLanguage, error,
) {
	call := service.I18nLanguages.List(parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetI18nLanguage, err)
	}

	return res.Items, nil
}

func (i *i18nLanguage) List(
	parts []string, output string, jsonpath string, writer io.Writer,
) error {
	i18nLanguages, err := i.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(i18nLanguages, jsonpath, writer)
	case "yaml":
		utils.PrintYAML(i18nLanguages, jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Hl", "Name"})
		for _, lang := range i18nLanguages {
			tb.AppendRow(table.Row{lang.Id, lang.Snippet.Hl, lang.Snippet.Name})
		}
	}
	return nil
}

func WithHl(hl string) Option {
	return func(i *i18nLanguage) {
		i.Hl = hl
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *i18nLanguage) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
