// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetMembershipsLevel = errors.New("failed to get memberships level")
)

type MembershipsLevel struct {
	*common.Fields
}

type IMembershipsLevel[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*MembershipsLevel)

func NewMembershipsLevel(opts ...Option) IMembershipsLevel[youtube.MembershipsLevel] {
	m := &MembershipsLevel{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MembershipsLevel) GetFields() *common.Fields {
	return m.Fields
}

func (m *MembershipsLevel) Get() ([]*youtube.MembershipsLevel, error) {
	if err := m.EnsureService(); err != nil {
		return nil, err
	}
	call := m.Service.MembershipsLevels.List(m.Parts)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetMembershipsLevel, err)
	}

	return res.Items, nil
}

func (m *MembershipsLevel) List(writer io.Writer) error {
	levels, err := m.Get()
	if err != nil {
		return err
	}

	common.PrintList(
		m.Output, levels, writer, table.Row{"ID", "Display Name"},
		func(ml *youtube.MembershipsLevel) table.Row {
			return table.Row{ml.Id, ml.Snippet.LevelDetails.DisplayName}
		},
	)
	return nil
}

var (
	WithParts   = common.WithParts[*MembershipsLevel]
	WithOutput  = common.WithOutput[*MembershipsLevel]
	WithService = common.WithService[*MembershipsLevel]
)
