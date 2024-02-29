package yutuber

import (
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

type PlaylistItem struct {
	id         string
	title      string
	desc       string
	playlistId string
	channelId  string
	privacy    string
}

type PlaylistItemService interface {
	List([]string, string)
	get([]string) []*youtube.PlaylistItem
}

type PlaylistItemOption func(*PlaylistItem)

func NewPlaylistItem(opts ...PlaylistItemOption) *PlaylistItem {
	p := &PlaylistItem{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (pi *PlaylistItem) List(parts []string, output string) {
	playlistItems := pi.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(playlistItems)
	case "yaml":
		utils.PrintYAML(playlistItems)
	default:
		fmt.Println("ID\tTitle")
		for _, playlistItem := range playlistItems {
			fmt.Printf("%s\t%s\n", playlistItem.Id, playlistItem.Snippet.Title)
		}
	}
}

func (pi *PlaylistItem) get(parts []string) []*youtube.PlaylistItem {
	call := service.PlaylistItems.List(parts)
	if pi.id != "" {
		call = call.Id(pi.id)
	} else if pi.playlistId != "" {
		call = call.PlaylistId(pi.playlistId)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist items: %v", err)
	}

	return response.Items
}

func WithPlaylistItemId(id string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.id = id
	}
}

func WithPlaylistItemTitle(title string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.title = title
	}
}

func WithPlaylistItemDesc(desc string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.desc = desc
	}
}

func WithPlaylistItemPlaylistId(playlistId string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.playlistId = playlistId
	}
}

func WithPlaylistItemChannelId(channelId string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.channelId = channelId
	}
}

func WithPlaylistItemPrivacy(privacy string) PlaylistItemOption {
	return func(p *PlaylistItem) {
		p.privacy = privacy
	}
}
