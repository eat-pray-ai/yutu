package yutuber

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylistItem    = fmt.Errorf("failed to get playlist item")
	errUpdatePlaylistItem = fmt.Errorf("failed to update playlist item")
	errInsertPlaylistItem = fmt.Errorf("failed to insert playlist item")
)

type playlistItem struct {
	id         string
	title      string
	desc       string
	videoId    string
	playlistId string
	channelId  string
	privacy    string
}

type PlaylistItem interface {
	List([]string, string)
	Insert()
	Update()
	get([]string) []*youtube.PlaylistItem
}

type PlaylistItemOption func(*playlistItem)

func NewPlaylistItem(opts ...PlaylistItemOption) PlaylistItem {
	p := &playlistItem{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (pi *playlistItem) List(parts []string, output string) {
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

func (pi *playlistItem) Insert() {
	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			Title:       pi.title,
			Description: pi.desc,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: pi.videoId,
			},
			PlaylistId: pi.playlistId,
			ChannelId:  pi.channelId,
		},
		Status: &youtube.PlaylistItemStatus{
			PrivacyStatus: pi.privacy,
		},
	}

	call := service.PlaylistItems.Insert([]string{"snippet", "status"}, playlistItem)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertPlaylistItem, err), pi.videoId)
	}
	fmt.Println("PlaylistItem inserted:")
	utils.PrintJSON(res)
}

func (pi *playlistItem) Update() {
	playlistItem := pi.get([]string{"id", "snippet", "status"})[0]
	if pi.title != "" {
		playlistItem.Snippet.Title = pi.title
	}
	if pi.desc != "" {
		playlistItem.Snippet.Description = pi.desc
	}
	if pi.privacy != "" {
		playlistItem.Status.PrivacyStatus = pi.privacy
	}

	call := service.PlaylistItems.Update([]string{"snippet", "status"}, playlistItem)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdatePlaylistItem, err), pi.id)
	}
	fmt.Println("PlaylistItem updated:")
	utils.PrintJSON(res)
}

func (pi *playlistItem) get(parts []string) []*youtube.PlaylistItem {
	call := service.PlaylistItems.List(parts)
	if pi.id != "" {
		call = call.Id(pi.id)
	} else if pi.playlistId != "" {
		call = call.PlaylistId(pi.playlistId)
	}
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetPlaylistItem, err), pi.id)
	}

	return response.Items
}

func WithPlaylistItemId(id string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.id = id
	}
}

func WithPlaylistItemTitle(title string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.title = title
	}
}

func WithPlaylistItemDesc(desc string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.desc = desc
	}
}

func WithPlaylistItemVideoId(videoId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.videoId = videoId
	}
}

func WithPlaylistItemPlaylistId(playlistId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.playlistId = playlistId
	}
}

func WithPlaylistItemChannelId(channelId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.channelId = channelId
	}
}

func WithPlaylistItemPrivacy(privacy string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.privacy = privacy
	}
}
