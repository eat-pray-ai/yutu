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
	service                      *youtube.Service
	errGetVideoAbuseReportReason = errors.New("failed to get video abuse report reason")
)

type videoAbuseReportReason struct {
	Hl string `yaml:"hl" json:"hl"`
}

type VideoAbuseReportReason[T any] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
}

type Option func(*videoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) VideoAbuseReportReason[youtube.VideoAbuseReportReason] {
	va := &videoAbuseReportReason{}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (va *videoAbuseReportReason) Get(parts []string) (
	[]*youtube.VideoAbuseReportReason, error,
) {
	call := service.VideoAbuseReportReasons.List(parts)
	if va.Hl != "" {
		call = call.Hl(va.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoAbuseReportReason, err)
	}

	return res.Items, nil
}

func (va *videoAbuseReportReason) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	videoAbuseReportReasons, err := va.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(videoAbuseReportReasons, jpath, writer)
	case "yaml":
		utils.PrintYAML(videoAbuseReportReasons, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Label"})
		for _, reason := range videoAbuseReportReasons {
			tb.AppendRow(table.Row{reason.Id, reason.Snippet.Label})
		}
	}
	return nil
}

func WithHL(hl string) Option {
	return func(vc *videoAbuseReportReason) {
		vc.Hl = hl
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *videoAbuseReportReason) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
