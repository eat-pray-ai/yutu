// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

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
	service                *youtube.Service
	errGetMembershipsLevel = errors.New("failed to get memberships level")
)

type membershipsLevel struct{}

type MembershipsLevel[T any] interface {
	List([]string, string, string, io.Writer) error
	Get([]string) ([]*T, error)
}

type Option func(*membershipsLevel)

func NewMembershipsLevel(opts ...Option) MembershipsLevel[youtube.MembershipsLevel] {
	m := &membershipsLevel{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *membershipsLevel) Get(parts []string) (
	[]*youtube.MembershipsLevel, error,
) {
	call := service.MembershipsLevels.List(parts)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetMembershipsLevel, err)
	}

	return res.Items, nil
}

func (m *membershipsLevel) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	membershipsLevels, err := m.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(membershipsLevels, jpath, writer)
	case "yaml":
		utils.PrintYAML(membershipsLevels, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Display Name"})
		for _, ml := range membershipsLevels {
			tb.AppendRow(table.Row{ml.Id, ml.Snippet.LevelDetails.DisplayName})
		}
	}
	return nil
}

func WithService(svc *youtube.Service) Option {
	return func(_ *membershipsLevel) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
