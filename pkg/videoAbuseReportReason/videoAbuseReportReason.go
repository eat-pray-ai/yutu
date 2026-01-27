// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

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
	errGetVideoAbuseReportReason = errors.New("failed to get video abuse report reason")
)

type VideoAbuseReportReason struct {
	*pkg.DefaultFields
	Hl string `yaml:"hl" json:"hl"`
}

type IVideoAbuseReportReason[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(*VideoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) IVideoAbuseReportReason[youtube.VideoAbuseReportReason] {
	va := &VideoAbuseReportReason{DefaultFields: &pkg.DefaultFields{}}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (va *VideoAbuseReportReason) GetDefaultFields() *pkg.DefaultFields {
	return va.DefaultFields
}

func (va *VideoAbuseReportReason) preRun() {
	if va.Service == nil {
		va.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (va *VideoAbuseReportReason) Get() (
	[]*youtube.VideoAbuseReportReason, error,
) {
	va.preRun()
	call := va.Service.VideoAbuseReportReasons.List(va.Parts)
	if va.Hl != "" {
		call = call.Hl(va.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoAbuseReportReason, err)
	}

	return res.Items, nil
}

func (va *VideoAbuseReportReason) List(writer io.Writer) error {
	videoAbuseReportReasons, err := va.Get()
	if err != nil {
		return err
	}

	switch va.Output {
	case "json":
		utils.PrintJSON(videoAbuseReportReasons, va.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(videoAbuseReportReasons, va.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Label"})
		for _, reason := range videoAbuseReportReasons {
			tb.AppendRow(table.Row{reason.Id, reason.Snippet.Label})
		}
	}
	return nil
}

func WithHL(hl string) Option {
	return func(va *VideoAbuseReportReason) {
		va.Hl = hl
	}
}

var (
	WithParts    = pkg.WithParts[*VideoAbuseReportReason]
	WithOutput   = pkg.WithOutput[*VideoAbuseReportReason]
	WithJsonpath = pkg.WithJsonpath[*VideoAbuseReportReason]
	WithService  = pkg.WithService[*VideoAbuseReportReason]
)
