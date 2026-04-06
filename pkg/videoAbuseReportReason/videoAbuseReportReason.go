// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideoAbuseReportReason = errors.New("failed to get video abuse report reason")
)

type VideoAbuseReportReason struct {
	*common.Fields
}

type IVideoAbuseReportReason[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*VideoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) IVideoAbuseReportReason[youtube.VideoAbuseReportReason] {
	va := &VideoAbuseReportReason{Fields: &common.Fields{}}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (va *VideoAbuseReportReason) Get() (
	[]*youtube.VideoAbuseReportReason, error,
) {
	if err := va.EnsureService(); err != nil {
		return nil, err
	}
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
	reasons, err := va.Get()
	if err != nil {
		return err
	}

	common.PrintList(va.Output, reasons, writer, table.Row{"ID", "Label"}, func(r *youtube.VideoAbuseReportReason) table.Row {
		return table.Row{r.Id, r.Snippet.Label}
	})
	return nil
}

var (
	WithHL      = common.WithHl[*VideoAbuseReportReason]
	WithParts   = common.WithParts[*VideoAbuseReportReason]
	WithOutput  = common.WithOutput[*VideoAbuseReportReason]
	WithService = common.WithService[*VideoAbuseReportReason]
)
