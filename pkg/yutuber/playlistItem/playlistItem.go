package playlistItem

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service               *youtube.Service
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

type Option func(*playlistItem)

func NewPlaylistItem(opts ...Option) PlaylistItem {
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

func WithId(id string) Option {
	return func(p *playlistItem) {
		p.id = id
	}
}

func WithTitle(title string) Option {
	return func(p *playlistItem) {
		p.title = title
	}
}

func WithDescription(description string) Option {
	return func(p *playlistItem) {
		p.description = description
	}
}

func WithKind(kind string) Option {
	return func(p *playlistItem) {
		p.kind = kind
	}
}

func WithKVideoId(kVideoId string) Option {
	return func(p *playlistItem) {
		p.kVideoId = kVideoId
	}
}

func WithKChannelId(kChannelId string) Option {
	return func(p *playlistItem) {
		p.kChannelId = kChannelId
	}
}

func WithKPlaylistId(kPlaylistId string) Option {
	return func(p *playlistItem) {
		p.kPlaylistId = kPlaylistId
	}
}

func WithVideoId(videoId string) Option {
	return func(p *playlistItem) {
		p.videoId = videoId
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(p *playlistItem) {
		p.playlistId = playlistId
	}
}

func WithChannelId(channelId string) Option {
	return func(p *playlistItem) {
		p.channelId = channelId
	}
}

func WithPrivacy(privacy string) Option {
	return func(p *playlistItem) {
		p.privacy = privacy
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(p *playlistItem) {
		p.maxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(p *playlistItem) {
		p.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}
