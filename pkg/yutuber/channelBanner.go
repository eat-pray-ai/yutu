package yutuber

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

type ChannelBannerOption func(banner *channelBanner)

func NewChannelBanner(opts ...ChannelBannerOption) ChannelBanner {
	service = auth.NewY2BService()
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
	utils.PrintJSON(res)
}

func WithChannelBannerFile(file string) ChannelBannerOption {
	return func(cb *channelBanner) {
		cb.file = file
	}
}

func WithChannelBannerOnBehalfOfContentOwner(onBehalfOfContentOwner string) ChannelBannerOption {
	return func(cb *channelBanner) {
		cb.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithChannelBannerOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) ChannelBannerOption {
	return func(cb *channelBanner) {
		cb.onBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}
