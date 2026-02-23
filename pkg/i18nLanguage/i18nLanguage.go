// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetI18nLanguage = errors.New("failed to get i18n language")
)

type I18nLanguage struct {
	*common.Fields
	Hl string `yaml:"hl" json:"hl,omitempty"`
}

type II18nLanguage[T youtube.I18nLanguage] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*I18nLanguage)

func NewI18nLanguage(opts ...Option) II18nLanguage[youtube.I18nLanguage] {
	i := &I18nLanguage{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *I18nLanguage) Get() (
	[]*youtube.I18nLanguage, error,
) {
	i.EnsureService()
	call := i.Service.I18nLanguages.List(i.Parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetI18nLanguage, err)
	}

	return res.Items, nil
}

func (i *I18nLanguage) List(writer io.Writer) error {
	i18nLanguages, err := i.Get()
	if err != nil {
		return err
	}

	switch i.Output {
	case "json":
		utils.PrintJSON(i18nLanguages, i.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(i18nLanguages, i.Jsonpath, writer)
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
	return func(i *I18nLanguage) {
		i.Hl = hl
	}
}

var (
	WithParts    = common.WithParts[*I18nLanguage]
	WithOutput   = common.WithOutput[*I18nLanguage]
	WithJsonpath = common.WithJsonpath[*I18nLanguage]
	WithService  = common.WithService[*I18nLanguage]
)
