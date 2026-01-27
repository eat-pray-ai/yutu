// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service               *youtube.Service
	errGetPlaylistItem    = errors.New("failed to get playlist item")
	errUpdatePlaylistItem = errors.New("failed to update playlist item")
	errInsertPlaylistItem = errors.New("failed to insert playlist item")
	errDeletePlaylistItem = errors.New("failed to delete playlist item")
)

type PlaylistItem struct {
	*pkg.DefaultFields
	Ids         []string `yaml:"ids" json:"ids"`
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

type IPlaylistItem[T any] interface {
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
	Get() ([]*T, error)
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(*PlaylistItem)

func NewPlaylistItem(opts ...Option) IPlaylistItem[youtube.PlaylistItem] {
	p := &PlaylistItem{DefaultFields: &pkg.DefaultFields{}}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *PlaylistItem) GetDefaultFields() *pkg.DefaultFields {
	return p.DefaultFields
}

func (pi *PlaylistItem) preRun() {
	if pi.Service == nil {
		pi.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (pi *PlaylistItem) Get() ([]*youtube.PlaylistItem, error) {
	pi.preRun()
	call := pi.Service.PlaylistItems.List(pi.Parts)
	if len(pi.Ids) > 0 {
		call = call.Id(pi.Ids...)
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

	var items []*youtube.PlaylistItem
	pageToken := ""
	for pi.MaxResults > 0 {
		call = call.MaxResults(min(pi.MaxResults, pkg.PerPage))
		pi.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetPlaylistItem, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (pi *PlaylistItem) List(writer io.Writer) error {
	playlistItems, err := pi.Get()
	if err != nil && playlistItems == nil {
		return err
	}

	switch pi.Output {
	case "json":
		utils.PrintJSON(playlistItems, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(playlistItems, pi.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
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
	return err
}

func (pi *PlaylistItem) Insert(writer io.Writer) error {
	pi.preRun()
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

	call := pi.Service.PlaylistItems.Insert(
		[]string{"snippet", "status"}, playlistItem,
	)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertPlaylistItem, err)
	}

	switch pi.Output {
	case "json":
		utils.PrintJSON(res, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, pi.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist Item inserted: %s\n", res.Id)
	}
	return nil
}

func (pi *PlaylistItem) Update(writer io.Writer) error {
	pi.preRun()
	pi.Parts = []string{"id", "snippet", "status"}
	playlistItems, err := pi.Get()

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

	call := pi.Service.PlaylistItems.Update(
		[]string{"snippet", "status"}, playlistItem,
	)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdatePlaylistItem, err)
	}

	switch pi.Output {
	case "json":
		utils.PrintJSON(res, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, pi.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist Item updated: %s\n", res.Id)
	}
	return nil
}

func (pi *PlaylistItem) Delete(writer io.Writer) error {
	pi.preRun()
	for _, id := range pi.Ids {
		call := pi.Service.PlaylistItems.Delete(id)
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

func WithIds(ids []string) Option {
	return func(p *PlaylistItem) {
		p.Ids = ids
	}
}

func WithTitle(title string) Option {
	return func(p *PlaylistItem) {
		p.Title = title
	}
}

func WithDescription(description string) Option {
	return func(p *PlaylistItem) {
		p.Description = description
	}
}

func WithKind(kind string) Option {
	return func(p *PlaylistItem) {
		p.Kind = kind
	}
}

func WithKVideoId(kVideoId string) Option {
	return func(p *PlaylistItem) {
		p.KVideoId = kVideoId
	}
}

func WithKChannelId(kChannelId string) Option {
	return func(p *PlaylistItem) {
		p.KChannelId = kChannelId
	}
}

func WithKPlaylistId(kPlaylistId string) Option {
	return func(p *PlaylistItem) {
		p.KPlaylistId = kPlaylistId
	}
}

func WithVideoId(videoId string) Option {
	return func(p *PlaylistItem) {
		p.VideoId = videoId
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(p *PlaylistItem) {
		p.PlaylistId = playlistId
	}
}

func WithChannelId(channelId string) Option {
	return func(p *PlaylistItem) {
		p.ChannelId = channelId
	}
}

func WithPrivacy(privacy string) Option {
	return func(p *PlaylistItem) {
		p.Privacy = privacy
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(p *PlaylistItem) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		p.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(p *PlaylistItem) {
		p.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

var (
	WithParts    = pkg.WithParts[*PlaylistItem]
	WithOutput   = pkg.WithOutput[*PlaylistItem]
	WithJsonpath = pkg.WithJsonpath[*PlaylistItem]
	WithService  = pkg.WithService[*PlaylistItem]
)
