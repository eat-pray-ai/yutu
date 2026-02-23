// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

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
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type I18nRegion struct {
	*common.Fields
	Hl string `yaml:"hl" json:"hl,omitempty"`
}

type II18nRegion[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*I18nRegion)

func NewI18nRegion(opts ...Option) II18nRegion[youtube.I18nRegion] {
	i := &I18nRegion{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *I18nRegion) Get() ([]*youtube.I18nRegion, error) {
	i.EnsureService()
	call := i.Service.I18nRegions.List(i.Parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetI18nRegion, err)
	}

	return res.Items, nil
}

func (i *I18nRegion) List(writer io.Writer) error {
	i18nRegions, err := i.Get()
	if err != nil {
		return err
	}

	switch i.Output {
	case "json":
		utils.PrintJSON(i18nRegions, i.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(i18nRegions, i.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Gl", "Name"})
		for _, region := range i18nRegions {
			tb.AppendRow(table.Row{region.Id, region.Snippet.Gl, region.Snippet.Name})
		}
	}
	return nil
}

func WithHl(hl string) Option {
	return func(i *I18nRegion) {
		i.Hl = hl
	}
}

var (
	WithParts    = common.WithParts[*I18nRegion]
	WithOutput   = common.WithOutput[*I18nRegion]
	WithJsonpath = common.WithJsonpath[*I18nRegion]
	WithService  = common.WithService[*I18nRegion]
)
