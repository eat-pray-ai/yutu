// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	errInsertChannelBanner = errors.New("failed to insert channelBanner")
)

type ChannelBanner struct {
	*pkg.DefaultFields
	ChannelId string `yaml:"channel_id" json:"channel_id"`
	File      string `yaml:"file" json:"file"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type IChannelBanner interface {
	Insert(io.Writer) error
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(banner *ChannelBanner)

func NewChannelBanner(opts ...Option) IChannelBanner {
	cb := &ChannelBanner{DefaultFields: &pkg.DefaultFields{}}
	for _, opt := range opts {
		opt(cb)
	}
	return cb
}

func (cb *ChannelBanner) GetDefaultFields() *pkg.DefaultFields {
	return cb.DefaultFields
}

func (cb *ChannelBanner) preRun() {
	if cb.Service == nil {
		cb.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (cb *ChannelBanner) Insert(writer io.Writer) error {
	cb.preRun()
	file, err := pkg.Root.Open(cb.File)
	if err != nil {
		return errors.Join(errInsertChannelBanner, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	cbr := &youtube.ChannelBannerResource{}

	call := cb.Service.ChannelBanners.Insert(cbr).ChannelId(cb.ChannelId).Media(file)
	if cb.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cb.OnBehalfOfContentOwner)
	}
	if cb.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(cb.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertChannelBanner, err)
	}

	switch cb.Output {
	case "json":
		utils.PrintJSON(res, cb.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, cb.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "ChannelBanner inserted: %s\n", res.Url)
	}
	return nil
}

func WithChannelId(channelId string) Option {
	return func(cb *ChannelBanner) {
		cb.ChannelId = channelId
	}
}

func WithFile(file string) Option {
	return func(cb *ChannelBanner) {
		cb.File = file
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cb *ChannelBanner) {
		cb.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(cb *ChannelBanner) {
		cb.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

var (
	WithOutput   = pkg.WithOutput[*ChannelBanner]
	WithJsonpath = pkg.WithJsonpath[*ChannelBanner]
	WithService  = pkg.WithService[*ChannelBanner]
)
