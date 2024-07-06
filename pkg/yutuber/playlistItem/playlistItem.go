package playlistItem

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"log"

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
	ID          string `yaml:"id" json:"id"`
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Kind        string `yaml:"kind" json:"kind"`
	KVideoId    string `yaml:"k_video_id" json:"k_video_id"`
	KChannelId  string `yaml:"k_channel_id" json:"k_channel_id"`
	KPlaylistId string `yaml:"k_playlist_id" json:"k_playlist_id"`
	VideoId     string `yaml:"video_id" json:"video_id"`
	PlaylistId  string `yaml:"playlist_id" json:"playlist_id"`
	ChannelId   string `yaml:"channel_id" json:"channel_id"`
	Privacy     string `yaml:"privacy" json:"privacy"`
	MaxResults  int64  `yaml:"max_results" json:"max_results"`

	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
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

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (pi *playlistItem) get(parts []string) []*youtube.PlaylistItem {
	call := service.PlaylistItems.List(parts)
	if pi.ID != "" {
		call = call.Id(pi.ID)
	}
	if pi.PlaylistId != "" {
		call = call.PlaylistId(pi.PlaylistId)
	}
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.VideoId != "" {
		call = call.VideoId(pi.VideoId)
	}

	call = call.MaxResults(pi.MaxResults)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(pi)
		log.Fatalln(errors.Join(errGetPlaylistItem, err), pi.ID)
	}

	return res.Items
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
	switch pi.Kind {
	case "video":
		resourceId = &youtube.ResourceId{
			Kind:    "youtube#video",
			VideoId: pi.KVideoId,
		}
	case "channel":
		resourceId = &youtube.ResourceId{
			Kind:      "youtube#channel",
			ChannelId: pi.KChannelId,
		}
	case "playlist":
		resourceId = &youtube.ResourceId{
			Kind:       "youtube#playlist",
			PlaylistId: pi.KPlaylistId,
		}
	}

	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			Title:       pi.Title,
			Description: pi.Description,
			ResourceId:  resourceId,
			PlaylistId:  pi.PlaylistId,
			ChannelId:   pi.ChannelId,
		},
		Status: &youtube.PlaylistItemStatus{
			PrivacyStatus: pi.Privacy,
		},
	}

	call := service.PlaylistItems.Insert(
		[]string{"snippet", "status"}, playlistItem,
	)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(pi)
		log.Fatalln(errors.Join(errInsertPlaylistItem, err), pi.VideoId)
	}

	utils.PrintYAML(res)
}

func (pi *playlistItem) Update() {
	playlistItem := pi.get([]string{"id", "snippet", "status"})[0]
	if pi.Title != "" {
		playlistItem.Snippet.Title = pi.Title
	}
	if pi.Description != "" {
		playlistItem.Snippet.Description = pi.Description
	}
	if pi.Privacy != "" {
		playlistItem.Status.PrivacyStatus = pi.Privacy
	}

	call := service.PlaylistItems.Update(
		[]string{"snippet", "status"}, playlistItem,
	)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(pi)
		log.Fatalln(errors.Join(errUpdatePlaylistItem, err), pi.ID)
	}

	utils.PrintYAML(res)
}

func (pi *playlistItem) Delete() {
	call := service.PlaylistItems.Delete(pi.ID)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		utils.PrintJSON(pi)
		log.Fatalln(errors.Join(errDeletePlaylistItem, err), pi.ID)
	}

	fmt.Printf("Playlsit Item %s deleted", pi.ID)
}

func WithID(id string) Option {
	return func(p *playlistItem) {
		p.ID = id
	}
}

func WithTitle(title string) Option {
	return func(p *playlistItem) {
		p.Title = title
	}
}

func WithDescription(description string) Option {
	return func(p *playlistItem) {
		p.Description = description
	}
}

func WithKind(kind string) Option {
	return func(p *playlistItem) {
		p.Kind = kind
	}
}

func WithKVideoId(kVideoId string) Option {
	return func(p *playlistItem) {
		p.KVideoId = kVideoId
	}
}

func WithKChannelId(kChannelId string) Option {
	return func(p *playlistItem) {
		p.KChannelId = kChannelId
	}
}

func WithKPlaylistId(kPlaylistId string) Option {
	return func(p *playlistItem) {
		p.KPlaylistId = kPlaylistId
	}
}

func WithVideoId(videoId string) Option {
	return func(p *playlistItem) {
		p.VideoId = videoId
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(p *playlistItem) {
		p.PlaylistId = playlistId
	}
}

func WithChannelId(channelId string) Option {
	return func(p *playlistItem) {
		p.ChannelId = channelId
	}
}

func WithPrivacy(privacy string) Option {
	return func(p *playlistItem) {
		p.Privacy = privacy
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(p *playlistItem) {
		p.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(p *playlistItem) {
		p.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithService(svc *youtube.Service) Option {
	return func(p *playlistItem) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
