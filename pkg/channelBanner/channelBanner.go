package channelBanner

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"io"
	"os"
)

var (
	service                *youtube.Service
	errInsertChannelBanner = errors.New("failed to insert channelBanner")
)

type channelBanner struct {
	ChannelId string `yaml:"channel_id" json:"channel_id"`
	File      string `yaml:"file" json:"file"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type ChannelBanner interface {
	Insert(string, string, io.Writer) error
}

type Option func(banner *channelBanner)

func NewChannelBanner(opts ...Option) ChannelBanner {
	cb := &channelBanner{}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (cb *channelBanner) Insert(
	output string, jpath string, writer io.Writer,
) error {
	file, err := os.Open(cb.File)
	if err != nil {
		return errors.Join(errInsertChannelBanner, err)
	}
	defer file.Close()
	cbr := &youtube.ChannelBannerResource{}

	call := service.ChannelBanners.Insert(cbr).ChannelId(cb.ChannelId).Media(file)
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

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "ChannelBanner inserted: %s\n", res.Url)
	}
	return nil
}

func WithChannelId(channelId string) Option {
	return func(cb *channelBanner) {
		cb.ChannelId = channelId
	}
}

func WithFile(file string) Option {
	return func(cb *channelBanner) {
		cb.File = file
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cb *channelBanner) {
		cb.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(cb *channelBanner) {
		cb.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *channelBanner) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
