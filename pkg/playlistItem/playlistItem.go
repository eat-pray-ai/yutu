package playlistItem

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service               *youtube.Service
	errGetPlaylistItem    = errors.New("failed to get playlist item")
	errUpdatePlaylistItem = errors.New("failed to update playlist item")
	errInsertPlaylistItem = errors.New("failed to insert playlist item")
	errDeletePlaylistItem = errors.New("failed to delete playlist item")
)

type playlistItem struct {
	IDs         []string `yaml:"ids" json:"ids"`
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Kind        string   `yaml:"kind" json:"kind"`
	KVideoId    string   `yaml:"k_video_id" json:"k_video_id"`
	KChannelId  string   `yaml:"k_channel_id" json:"k_channel_id"`
	KPlaylistId string   `yaml:"k_playlist_id" json:"k_playlist_id"`
	VideoId     string   `yaml:"video_id" json:"video_id"`
	PlaylistId  string   `yaml:"playlist_id" json:"playlist_id"`
	ChannelId   string   `yaml:"channel_id" json:"channel_id"`
	Privacy     string   `yaml:"privacy" json:"privacy"`
	MaxResults  int64    `yaml:"max_results" json:"max_results"`

	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
}

type PlaylistItem interface {
	List([]string, string, io.Writer) error
	Insert(string, io.Writer) error
	Update(string, io.Writer) error
	Delete(io.Writer) error
	Get([]string) ([]*youtube.PlaylistItem, error)
}

type Option func(*playlistItem)

func NewPlaylistItem(opts ...Option) PlaylistItem {
	p := &playlistItem{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (pi *playlistItem) Get(parts []string) ([]*youtube.PlaylistItem, error) {
	call := service.PlaylistItems.List(parts)
	if len(pi.IDs) > 0 {
		call = call.Id(pi.IDs...)
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
	if pi.MaxResults <= 0 {
		pi.MaxResults = 1
	}
	call = call.MaxResults(pi.MaxResults)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetPlaylistItem, err)
	}

	return res.Items, nil
}

func (pi *playlistItem) List(
	parts []string, output string, writer io.Writer,
) error {
	playlistItems, err := pi.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(playlistItems, writer)
	case "yaml":
		utils.PrintYAML(playlistItems, writer)
	default:
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Title", "Kind", "Resource ID"})
		for _, item := range playlistItems {
			var resourceId string
			switch item.Snippet.ResourceId.Kind {
			case "youtube#video":
				resourceId = item.Snippet.ResourceId.VideoId
			case "youtube#channel":
				resourceId = item.Snippet.ResourceId.ChannelId
			case "youtube#playlist":
				resourceId = item.Snippet.ResourceId.PlaylistId
			}
			tb.AppendRow(
				table.Row{
					item.Id, item.Snippet.Title, item.Snippet.ResourceId.Kind, resourceId,
				},
			)
		}
	}
	return nil
}

func (pi *playlistItem) Insert(output string, writer io.Writer) error {
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
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertPlaylistItem, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, writer)
	case "yaml":
		utils.PrintYAML(res, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist Item inserted: %s\n", res.Id)
	}
	return nil
}

func (pi *playlistItem) Update(output string, writer io.Writer) error {
	playlistItems, err := pi.Get([]string{"id", "snippet", "status"})
	if err != nil {
		return errors.Join(errUpdatePlaylistItem, err)
	}
	if len(playlistItems) == 0 {
		return errGetPlaylistItem
	}

	playlistItem := playlistItems[0]
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
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdatePlaylistItem, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, writer)
	case "yaml":
		utils.PrintYAML(res, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist Item updated: %s\n", res.Id)
	}
	return nil
}

func (pi *playlistItem) Delete(writer io.Writer) error {
	for _, id := range pi.IDs {
		call := service.PlaylistItems.Delete(id)
		if pi.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeletePlaylistItem, err)
		}

		_, _ = fmt.Fprintf(writer, "Playlsit Item %s deleted", id)
	}
	return nil
}

func WithIDs(ids []string) Option {
	return func(p *playlistItem) {
		p.IDs = ids
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
	return func(_ *playlistItem) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
