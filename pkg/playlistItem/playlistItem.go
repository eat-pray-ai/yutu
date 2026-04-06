// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylistItem    = errors.New("failed to get playlist item")
	errUpdatePlaylistItem = errors.New("failed to update playlist item")
	errInsertPlaylistItem = errors.New("failed to insert playlist item")
	errDeletePlaylistItem = errors.New("failed to delete playlist item")
)

type PlaylistItem struct {
	*common.Fields
	Title       string `yaml:"title" json:"title,omitempty"`
	Description string `yaml:"description" json:"description,omitempty"`
	Kind        string `yaml:"kind" json:"kind,omitempty"`
	KVideoId    string `yaml:"k_video_id" json:"k_video_id,omitempty"`
	KChannelId  string `yaml:"k_channel_id" json:"k_channel_id,omitempty"`
	KPlaylistId string `yaml:"k_playlist_id" json:"k_playlist_id,omitempty"`
	VideoId     string `yaml:"video_id" json:"video_id,omitempty"`
	PlaylistId  string `yaml:"playlist_id" json:"playlist_id,omitempty"`
	Privacy     string `yaml:"privacy" json:"privacy,omitempty"`
}

type IPlaylistItem[T any] interface {
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*PlaylistItem)

func NewPlaylistItem(opts ...Option) IPlaylistItem[youtube.PlaylistItem] {
	p := &PlaylistItem{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (pi *PlaylistItem) Get() ([]*youtube.PlaylistItem, error) {
	if err := pi.EnsureService(); err != nil {
		return nil, err
	}
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

	return common.Paginate(pi.Fields, call, func(r *youtube.PlaylistItemListResponse) ([]*youtube.PlaylistItem, string) {
		return r.Items, r.NextPageToken
	}, errGetPlaylistItem)
}

func (pi *PlaylistItem) List(writer io.Writer) error {
	playlistItems, err := pi.Get()
	if err != nil && playlistItems == nil {
		return err
	}

	common.PrintList(pi.Output, playlistItems, writer, table.Row{"ID", "Title", "Kind", "Resource ID"}, func(item *youtube.PlaylistItem) table.Row {
		title := ""
		kind := ""
		resourceId := ""
		if item.Snippet != nil {
			title = item.Snippet.Title
			if item.Snippet.ResourceId != nil {
				kind = item.Snippet.ResourceId.Kind
				switch kind {
				case "youtube#video":
					resourceId = item.Snippet.ResourceId.VideoId
				case "youtube#channel":
					resourceId = item.Snippet.ResourceId.ChannelId
				case "youtube#playlist":
					resourceId = item.Snippet.ResourceId.PlaylistId
				}
			}
		}
		return table.Row{item.Id, title, kind, resourceId}
	})
	return err
}

func (pi *PlaylistItem) Insert(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
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

	common.PrintResult(pi.Output, res, writer, "Playlist Item inserted: %s\n", res.Id)
	return nil
}

func (pi *PlaylistItem) Update(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
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

	common.PrintResult(pi.Output, res, writer, "Playlist Item updated: %s\n", res.Id)
	return nil
}

func (pi *PlaylistItem) Delete(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
	for _, id := range pi.Ids {
		call := pi.Service.PlaylistItems.Delete(id)
		if pi.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeletePlaylistItem, err)
		}

		_, _ = fmt.Fprintf(writer, "Playlist Item %s deleted\n", id)
	}
	return nil
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

func WithPrivacy(privacy string) Option {
	return func(p *PlaylistItem) {
		p.Privacy = privacy
	}
}

var (
	WithParts      = common.WithParts[*PlaylistItem]
	WithOutput     = common.WithOutput[*PlaylistItem]
	WithService    = common.WithService[*PlaylistItem]
	WithIds        = common.WithIds[*PlaylistItem]
	WithMaxResults = common.WithMaxResults[*PlaylistItem]
	WithChannelId  = common.WithChannelId[*PlaylistItem]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*PlaylistItem]
)
