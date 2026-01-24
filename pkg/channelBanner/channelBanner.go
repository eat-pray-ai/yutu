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
	ChannelId string `yaml:"channel_id" json:"channel_id"`
	File      string `yaml:"file" json:"file"`
	Output    string `yaml:"output" json:"output"`
	Jsonpath  string `yaml:"jsonpath" json:"jsonpath"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`

	service *youtube.Service
}

type IChannelBanner interface {
	Insert(io.Writer) error
	preRun()
}

type Option func(banner *ChannelBanner)

func NewChannelBanner(opts ...Option) IChannelBanner {
	cb := &ChannelBanner{}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (cb *ChannelBanner) preRun() {
	if cb.service == nil {
		cb.service = auth.NewY2BService(
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

	call := cb.service.ChannelBanners.Insert(cbr).ChannelId(cb.ChannelId).Media(file)
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

func WithOutput(output string) Option {
	return func(cb *ChannelBanner) {
		cb.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(cb *ChannelBanner) {
		cb.Jsonpath = jsonpath
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

func WithService(svc *youtube.Service) Option {
	return func(cb *ChannelBanner) {
		cb.service = svc
	}
}
