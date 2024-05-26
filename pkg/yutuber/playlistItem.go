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
	errGetPlaylistItem    = errors.New("failed to get playlist item")
	errUpdatePlaylistItem = errors.New("failed to update playlist item")
	errInsertPlaylistItem = errors.New("failed to insert playlist item")
	errDeletePlaylistItem = errors.New("failed to delete playlist item")
)

type playlistItem struct {
	id          string
	title       string
	description string
	kind        string
	kVideoId    string
	kChannelId  string
	kPlaylistId string
	videoId     string
	playlistId  string
	channelId   string
	privacy     string
	maxResults  int64

	onBehalfOfContentOwner string
}

type PlaylistItem interface {
	List([]string, string)
	Insert()
	Update()
	Delete()
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

func (pi *playlistItem) get(parts []string) []*youtube.PlaylistItem {
	call := service.PlaylistItems.List(parts)
	if pi.id != "" {
		call = call.Id(pi.id)
	}
	if pi.playlistId != "" {
		call = call.PlaylistId(pi.playlistId)
	}
	if pi.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.onBehalfOfContentOwner)
	}
	if pi.videoId != "" {
		call = call.VideoId(pi.videoId)
	}

	call = call.MaxResults(pi.maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetPlaylistItem, err), pi.id)
	}

	return response.Items
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
	var resourceId *youtube.ResourceId
	switch pi.kind {
	case "video":
		resourceId = &youtube.ResourceId{
			Kind:    "youtube#video",
			VideoId: pi.kVideoId,
		}
	case "channel":
		resourceId = &youtube.ResourceId{
			Kind:      "youtube#channel",
			ChannelId: pi.kChannelId,
		}
	case "playlist":
		resourceId = &youtube.ResourceId{
			Kind:       "youtube#playlist",
			PlaylistId: pi.kPlaylistId,
		}
	}

	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			Title:       pi.title,
			Description: pi.description,
			ResourceId:  resourceId,
			PlaylistId:  pi.playlistId,
			ChannelId:   pi.channelId,
		},
		Status: &youtube.PlaylistItemStatus{
			PrivacyStatus: pi.privacy,
		},
	}

	call := service.PlaylistItems.Insert(
		[]string{"snippet", "status"}, playlistItem,
	)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertPlaylistItem, err), pi.videoId)
	}
	fmt.Println("PlaylistItem inserted:")
	utils.PrintYAML(res)
}

func (pi *playlistItem) Update() {
	playlistItem := pi.get([]string{"id", "snippet", "status"})[0]
	if pi.title != "" {
		playlistItem.Snippet.Title = pi.title
	}
	if pi.description != "" {
		playlistItem.Snippet.Description = pi.description
	}
	if pi.privacy != "" {
		playlistItem.Status.PrivacyStatus = pi.privacy
	}

	call := service.PlaylistItems.Update(
		[]string{"snippet", "status"}, playlistItem,
	)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdatePlaylistItem, err), pi.id)
	}
	fmt.Println("PlaylistItem updated:")
	utils.PrintYAML(res)
}

func (pi *playlistItem) Delete() {
	call := service.PlaylistItems.Delete(pi.id)
	if pi.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeletePlaylistItem, err), pi.id)
	}

	fmt.Printf("Playlsit Item %s deleted", pi.id)
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

func WithPlaylistItemDescription(description string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.description = description
	}
}

func WithPlaylistItemKind(kind string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.kind = kind
	}
}

func WithPlaylistItemKVideoId(kVideoId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.kVideoId = kVideoId
	}
}

func WithPlaylistItemKChannelId(kChannelId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.kChannelId = kChannelId
	}
}

func WithPlaylistItemKPlaylistId(kPlaylistId string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.kPlaylistId = kPlaylistId
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

func WithPlaylistItemMaxResults(maxResults int64) PlaylistItemOption {
	return func(p *playlistItem) {
		p.maxResults = maxResults
	}
}

func WithPlaylistItemOnBehalfOfContentOwner(onBehalfOfContentOwner string) PlaylistItemOption {
	return func(p *playlistItem) {
		p.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}
