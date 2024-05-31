package channelBanner

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service                *youtube.Service
	errOpenFile            = errors.New("failed to open file")
	errInsertChannelBanner = errors.New("failed to insert channelBanner")
)

type channelBanner struct {
	file string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
}

type ChannelBanner interface {
	Insert()
}

type Option func(banner *channelBanner)

func NewChannelBanner(opts ...Option) ChannelBanner {
	cb := &channelBanner{}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (cb *channelBanner) Insert() {
	file, err := os.Open(cb.file)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err), cb.file)
	}
	defer file.Close()
	cbr := &youtube.ChannelBannerResource{}

	call := service.ChannelBanners.Insert(cbr).Media(file)
	if cb.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cb.onBehalfOfContentOwner)
	}
	if cb.onBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(cb.onBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertChannelBanner, err))
	}

	fmt.Println("ChannelBanner inserted:")
	utils.PrintYAML(res)
}

func WithFile(file string) Option {
	return func(cb *channelBanner) {
		cb.file = file
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cb *channelBanner) {
		cb.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(cb *channelBanner) {
		cb.onBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithService() Option {
	return func(cb *channelBanner) {
		service = auth.NewY2BService()
	}
}
