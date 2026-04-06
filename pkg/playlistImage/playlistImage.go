// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylistImage    = errors.New("failed to get playlist image")
	errInsertPlaylistImage = errors.New("failed to insert playlist image")
	errUpdatePlaylistImage = errors.New("failed to update playlist image")
	errDeletePlaylistImage = errors.New("failed to delete playlist image")
)

type PlaylistImage struct {
	*common.Fields
	Height     int64  `yaml:"height" json:"height,omitempty"`
	PlaylistId string `yaml:"playlist_id" json:"playlist_id,omitempty"`
	Type       string `yaml:"type" json:"type,omitempty"`
	Width      int64  `yaml:"width" json:"width,omitempty"`
	File       string `yaml:"file" json:"file,omitempty"`
	Parent     string `yaml:"parent" json:"parent,omitempty"`

	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type IPlaylistImage[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*PlaylistImage)

func NewPlaylistImage(opts ...Option) IPlaylistImage[youtube.PlaylistImage] {
	pi := &PlaylistImage{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(pi)
	}
	return pi
}

func (pi *PlaylistImage) Get() ([]*youtube.PlaylistImage, error) {
	if err := pi.EnsureService(); err != nil {
		return nil, err
	}
	call := pi.Service.PlaylistImages.List()
	call = call.Part(pi.Parts...)
	if pi.Parent != "" {
		call = call.Parent(pi.Parent)
	}
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(pi.OnBehalfOfContentOwnerChannel)
	}

	return common.Paginate(pi.Fields, call, func(r *youtube.PlaylistImageListResponse) ([]*youtube.PlaylistImage, string) {
		return r.Items, r.NextPageToken
	}, errGetPlaylistImage)
}

func (pi *PlaylistImage) List(writer io.Writer) error {
	playlistImages, err := pi.Get()
	if err != nil && playlistImages == nil {
		return err
	}

	common.PrintList(pi.Output, playlistImages, writer, table.Row{"ID", "Kind", "Playlist ID", "Type"}, func(img *youtube.PlaylistImage) table.Row {
		return table.Row{img.Id, img.Kind, img.Snippet.PlaylistId, img.Snippet.Type}
	})
	return err
}

func (pi *PlaylistImage) Insert(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
	file, err := pkg.Root.Open(pi.File)
	if err != nil {
		return errors.Join(errInsertPlaylistImage, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	playlistImage := &youtube.PlaylistImage{
		Kind: "youtube#playlistImages",
		Snippet: &youtube.PlaylistImageSnippet{
			PlaylistId: pi.PlaylistId,
			Type:       pi.Type,
			Height:     pi.Height,
			Width:      pi.Width,
		},
	}

	call := pi.Service.PlaylistImages.Insert(playlistImage)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(pi.OnBehalfOfContentOwnerChannel)
	}
	call = call.Media(file)
	call = call.Part("kind", "snippet")
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertPlaylistImage, err)
	}

	common.PrintResult(pi.Output, res, writer, "PlaylistImage inserted: %s\n", res.Id)
	return nil
}

func (pi *PlaylistImage) Update(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
	pi.Parts = []string{"id", "kind", "snippet"}
	playlistImages, err := pi.Get()
	if err != nil {
		return errors.Join(errUpdatePlaylistImage, err)
	}
	if len(playlistImages) == 0 {
		return errGetPlaylistImage
	}

	playlistImage := playlistImages[0]
	if pi.PlaylistId != "" {
		playlistImage.Snippet.PlaylistId = pi.PlaylistId
	}
	if pi.Type != "" {
		playlistImage.Snippet.Type = pi.Type
	}
	if pi.Height != 0 {
		playlistImage.Snippet.Height = pi.Height
	}
	if pi.Width != 0 {
		playlistImage.Snippet.Width = pi.Width
	}

	call := pi.Service.PlaylistImages.Update(playlistImage)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.File != "" {
		file, err := pkg.Root.Open(pi.File)
		if err != nil {
			return errors.Join(errUpdatePlaylistImage, err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		call = call.Media(file)
	}
	call = call.Part("id", "kind", "snippet")

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdatePlaylistImage, err)
	}

	common.PrintResult(pi.Output, res, writer, "PlaylistImage updated: %s\n", res.Id)
	return nil
}

func (pi *PlaylistImage) Delete(writer io.Writer) error {
	if err := pi.EnsureService(); err != nil {
		return err
	}
	for _, id := range pi.Ids {
		call := pi.Service.PlaylistImages.Delete()
		call = call.Id(id)
		if pi.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeletePlaylistImage, err)
		}
		_, _ = fmt.Fprintf(writer, "PlaylistImage %s deleted\n", id)
	}
	return nil
}

func WithHeight(height int64) Option {
	return func(pi *PlaylistImage) {
		pi.Height = height
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(pi *PlaylistImage) {
		pi.PlaylistId = playlistId
	}
}

func WithType(t string) Option {
	return func(pi *PlaylistImage) {
		pi.Type = t
	}
}

func WithWidth(width int64) Option {
	return func(pi *PlaylistImage) {
		pi.Width = width
	}
}

func WithFile(file string) Option {
	return func(pi *PlaylistImage) {
		pi.File = file
	}
}

func WithParent(parent string) Option {
	return func(pi *PlaylistImage) {
		pi.Parent = parent
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(pi *PlaylistImage) {
		pi.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

var (
	WithParts      = common.WithParts[*PlaylistImage]
	WithOutput     = common.WithOutput[*PlaylistImage]
	WithService    = common.WithService[*PlaylistImage]
	WithIds        = common.WithIds[*PlaylistImage]
	WithMaxResults = common.WithMaxResults[*PlaylistImage]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*PlaylistImage]
)
