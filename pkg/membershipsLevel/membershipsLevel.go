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
	errGetMembershipsLevel = errors.New("failed to get memberships level")
)

type MembershipsLevel struct {
	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
	service  *youtube.Service
}

type IMembershipsLevel[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	preRun()
}

type Option func(*MembershipsLevel)

func NewMembershipsLevel(opts ...Option) IMembershipsLevel[youtube.MembershipsLevel] {
	m := &MembershipsLevel{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *MembershipsLevel) Get() (
	[]*youtube.MembershipsLevel, error,
) {
	m.preRun()
	call := m.service.MembershipsLevels.List(m.Parts)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetMembershipsLevel, err)
	}

	return res.Items, nil
}

func (m *MembershipsLevel) List(writer io.Writer) error {
	membershipsLevels, err := m.Get()
	if err != nil {
		return err
	}

	switch m.Output {
	case "json":
		utils.PrintJSON(membershipsLevels, m.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(membershipsLevels, m.Jsonpath, writer)
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

func (m *MembershipsLevel) preRun() {
	if m.service == nil {
		m.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func WithService(svc *youtube.Service) Option {
	return func(m *MembershipsLevel) {
		m.service = svc
	}
}

func WithParts(parts []string) Option {
	return func(m *MembershipsLevel) {
		m.Parts = parts
	}
}

func WithOutput(output string) Option {
	return func(m *MembershipsLevel) {
		m.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(m *MembershipsLevel) {
		m.Jsonpath = jsonpath
	}
}
