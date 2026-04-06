// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetI18nLanguage = errors.New("failed to get i18n language")
)

type I18nLanguage struct {
	*common.Fields
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
	if err := i.EnsureService(); err != nil {
		return nil, err
	}
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
	languages, err := i.Get()
	if err != nil {
		return err
	}

	common.PrintList(
		i.Output, languages, writer, table.Row{"ID", "Hl", "Name"},
		func(l *youtube.I18nLanguage) table.Row {
			return table.Row{l.Id, l.Snippet.Hl, l.Snippet.Name}
		},
	)
	return nil
}

var (
	WithHl      = common.WithHl[*I18nLanguage]
	WithParts   = common.WithParts[*I18nLanguage]
	WithOutput  = common.WithOutput[*I18nLanguage]
	WithService = common.WithService[*I18nLanguage]
)
