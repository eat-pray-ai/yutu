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
	Hl       string   `yaml:"hl" json:"hl"`
	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
	service  *youtube.Service
}

type IVideoAbuseReportReason[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	preRun()
}

type Option func(*VideoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) IVideoAbuseReportReason[youtube.VideoAbuseReportReason] {
	va := &VideoAbuseReportReason{}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (va *VideoAbuseReportReason) preRun() {
	if va.service == nil {
		va.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (va *VideoAbuseReportReason) Get() (
	[]*youtube.VideoAbuseReportReason, error,
) {
	va.preRun()
	call := va.service.VideoAbuseReportReasons.List(va.Parts)
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

func WithParts(parts []string) Option {
	return func(va *VideoAbuseReportReason) {
		va.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(va *VideoAbuseReportReason) {
		va.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(va *VideoAbuseReportReason) {
		va.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(va *VideoAbuseReportReason) {
		va.service = svc
	}
}
