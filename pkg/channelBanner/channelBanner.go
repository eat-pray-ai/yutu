// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

var (
	errInsertChannelBanner = errors.New("failed to insert channelBanner")
)

type ChannelBanner struct {
	common.Fields
	File string `yaml:"file" json:"file,omitempty"`

	OBOCOC string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type IChannelBanner interface {
	Insert(io.Writer) error
}

type Option func(banner *ChannelBanner)

func NewChannelBanner(opts ...Option) IChannelBanner {
	cb := &ChannelBanner{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

func (cb *ChannelBanner) Insert(writer io.Writer) error {
	if err := cb.EnsureService(); err != nil {
		return err
	}
	file, err := pkg.OpenFile(cb.File)
	if err != nil {
		return errors.Join(errInsertChannelBanner, err)
	}
	defer func() { _ = file.Close() }()
	cbr := &youtube.ChannelBannerResource{}

	call := cb.Service.ChannelBanners.Insert(cbr).ChannelId(cb.ChannelId).Media(file)
	if cb.OBOCO != "" {
		call = call.OnBehalfOfContentOwner(cb.OBOCO)
	}
	if cb.OBOCOC != "" {
		call = call.OnBehalfOfContentOwnerChannel(cb.OBOCOC)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertChannelBanner, err)
	}

	return common.PrintResult(
		cb.Output, res, writer, "ChannelBanner inserted: %s\n", res.Url,
	)
}

func WithFile(file string) Option {
	return func(cb *ChannelBanner) {
		cb.File = file
	}
}

func WithOBOCOC(obococ string) Option {
	return func(cb *ChannelBanner) {
		cb.OBOCOC = obococ
	}
}

var (
	WithChannelId = common.WithChannelId[*ChannelBanner]
	WithOutput    = common.WithOutput[*ChannelBanner]
	WithService   = common.WithService[*ChannelBanner]

	WithOBOCO = common.WithOBOCO[*ChannelBanner]
)
