package channelSection

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service                 *youtube.Service
	errGetChannelSection    = errors.New("failed to get channel section")
	errDeleteChannelSection = errors.New("failed to delete channel section")
)

type channelSection struct {
	id                     string
	channelId              string
	hl                     string
	mine                   *bool
	onBehalfOfContentOwner string
}

type ChannelSection interface {
	get(parts []string) []*youtube.ChannelSection
	List(parts []string, output string)
	// Update()
	// Insert()
	Delete()
}

type Option func(*channelSection)

func NewChannelSection(opts ...Option) ChannelSection {
	cs := &channelSection{}

	for _, opt := range opts {
		opt(cs)
	}
	return cs
}

func (cs *channelSection) get(parts []string) []*youtube.ChannelSection {
	call := service.ChannelSections.List(parts)
	if cs.id != "" {
		call = call.Id(cs.id)
	}
	if cs.channelId != "" {
		call = call.ChannelId(cs.channelId)
	}
	if cs.hl != "" {
		call = call.Hl(cs.hl)
	}
	if cs.mine != nil {
		call = call.Mine(*cs.mine)
	}
	if cs.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cs.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetChannelSection, err))
	}
	return res.Items
}

func (cs *channelSection) List(parts []string, output string) {
	channelSections := cs.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(channelSections)
	case "yaml":
		utils.PrintYAML(channelSections)
	default:
		fmt.Println("ID\tChannelID\tTitle")
		for _, channelSection := range channelSections {
			fmt.Printf(
				"%s\t%s\t%s\n", channelSection.Id,
				channelSection.Snippet.ChannelId, channelSection.Snippet.Title,
			)
		}
	}
}

func (cs *channelSection) Delete() {
	call := service.ChannelSections.Delete(cs.id)
	if cs.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cs.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeleteChannelSection, err))
	}

	fmt.Printf("Channel section %s deleted\n", cs.id)
}

func WithId(id string) Option {
	return func(cs *channelSection) {
		cs.id = id
	}
}

func WithChannelId(channelId string) Option {
	return func(cs *channelSection) {
		cs.channelId = channelId
	}
}

func WithHl(hl string) Option {
	return func(cs *channelSection) {
		cs.hl = hl
	}
}

func WithMine(mine bool, changed bool) Option {
	return func(cs *channelSection) {
		if changed {
			cs.mine = &mine
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cs *channelSection) {
		cs.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithService(svc *youtube.Service) Option {
	return func(cs *channelSection) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
