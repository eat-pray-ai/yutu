// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

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
	errGetVideoCategory = errors.New("failed to get video categoryId")
)

type VideoCategory struct {
	Ids        []string `yaml:"ids" json:"ids"`
	Hl         string   `yaml:"hl" json:"hl"`
	RegionCode string   `yaml:"region_code" json:"region_code"`
	*pkg.DefaultFields
}

type IVideoCategory[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(*VideoCategory)

func NewVideoCategory(opt ...Option) IVideoCategory[youtube.VideoCategory] {
	vc := &VideoCategory{DefaultFields: &pkg.DefaultFields{}}
	for _, o := range opt {
		o(vc)
	}
	return vc
}

func (vc *VideoCategory) GetDefaultFields() *pkg.DefaultFields {
	return vc.DefaultFields
}

func (vc *VideoCategory) preRun() {
	if vc.Service == nil {
		vc.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (vc *VideoCategory) Get() ([]*youtube.VideoCategory, error) {
	vc.preRun()
	call := vc.Service.VideoCategories.List(vc.Parts)
	if len(vc.Ids) > 0 {
		call = call.Id(vc.Ids...)
	}
	if vc.Hl != "" {
		call = call.Hl(vc.Hl)
	}
	if vc.RegionCode != "" {
		call = call.RegionCode(vc.RegionCode)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoCategory, err)
	}

	return res.Items, nil
}

func (vc *VideoCategory) List(writer io.Writer) error {
	videoCategories, err := vc.Get()
	if err != nil {
		return err
	}

	switch vc.Output {
	case "json":
		utils.PrintJSON(videoCategories, vc.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(videoCategories, vc.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Title", "Assignable"})
		for _, cat := range videoCategories {
			tb.AppendRow(table.Row{cat.Id, cat.Snippet.Title, cat.Snippet.Assignable})
		}
	}
	return nil
}

func WithIds(ids []string) Option {
	return func(vc *VideoCategory) {
		vc.Ids = ids
	}
}

func WithHl(hl string) Option {
	return func(vc *VideoCategory) {
		vc.Hl = hl
	}
}

func WithRegionCode(regionCode string) Option {
	return func(vc *VideoCategory) {
		vc.RegionCode = regionCode
	}
}

var (
	WithParts    = pkg.WithParts[*VideoCategory]
	WithOutput   = pkg.WithOutput[*VideoCategory]
	WithJsonpath = pkg.WithJsonpath[*VideoCategory]
	WithService  = pkg.WithService[*VideoCategory]
)
