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
	ID                     string `yaml:"id" json:"id"`
	ChannelId              string `yaml:"channel_id" json:"channel_id"`
	Hl                     string `yaml:"hl" json:"hl"`
	Mine                   *bool  `yaml:"mine" json:"mine"`
	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
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
	if cs.ID != "" {
		call = call.Id(cs.ID)
	}
	if cs.ChannelId != "" {
		call = call.ChannelId(cs.ChannelId)
	}
	if cs.Hl != "" {
		call = call.Hl(cs.Hl)
	}
	if cs.Mine != nil {
		call = call.Mine(*cs.Mine)
	}
	if cs.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cs.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(cs)
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
	call := service.ChannelSections.Delete(cs.ID)
	if cs.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(cs.OnBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		utils.PrintJSON(cs)
		log.Fatalln(errors.Join(errDeleteChannelSection, err))
	}

	fmt.Printf("Channel section %s deleted\n", cs.ID)
}

func WithID(id string) Option {
	return func(cs *channelSection) {
		cs.ID = id
	}
}

func WithChannelId(channelId string) Option {
	return func(cs *channelSection) {
		cs.ChannelId = channelId
	}
}

func WithHl(hl string) Option {
	return func(cs *channelSection) {
		cs.Hl = hl
	}
}

func WithMine(mine bool, changed bool) Option {
	return func(cs *channelSection) {
		if changed {
			cs.Mine = &mine
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(cs *channelSection) {
		cs.OnBehalfOfContentOwner = onBehalfOfContentOwner
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
