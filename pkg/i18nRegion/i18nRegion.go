// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

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
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type I18nRegion struct {
	Hl       string   `yaml:"hl" json:"hl"`
	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
	service  *youtube.Service
}

type II18nRegion[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	preRun()
}

type Option func(*I18nRegion)

func NewI18nRegion(opts ...Option) II18nRegion[youtube.I18nRegion] {
	i := &I18nRegion{}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func (i *I18nRegion) preRun() {
	if i.service == nil {
		i.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (i *I18nRegion) Get() ([]*youtube.I18nRegion, error) {
	i.preRun()
	call := i.service.I18nRegions.List(i.Parts)
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

func WithParts(parts []string) Option {
	return func(i *I18nRegion) {
		i.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(i *I18nRegion) {
		i.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(i *I18nRegion) {
		i.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(i *I18nRegion) {
		i.service = svc
	}
}
