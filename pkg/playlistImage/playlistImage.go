// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
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
	Ids        []string `yaml:"ids" json:"ids"`
	Height     int64    `yaml:"height" json:"height"`
	PlaylistId string   `yaml:"playlist_id" json:"playlist_id"`
	Type       string   `yaml:"type" json:"type"`
	Width      int64    `yaml:"width" json:"width"`
	File       string   `yaml:"file" json:"file"`

	Parent     string `yaml:"parent" json:"parent"`
	MaxResults int64  `yaml:"max_results" json:"max_results"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
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
	pi.EnsureService()
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

	var items []*youtube.PlaylistImage
	pageToken := ""
	for pi.MaxResults > 0 {
		call = call.MaxResults(min(pi.MaxResults, pkg.PerPage))
		pi.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetPlaylistImage, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (pi *PlaylistImage) List(writer io.Writer) error {
	playlistImages, err := pi.Get()
	if err != nil && playlistImages == nil {
		return err
	}

	switch pi.Output {
	case "json":
		utils.PrintJSON(playlistImages, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(playlistImages, pi.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Kind", "Playlist ID", "Type"})
		for _, img := range playlistImages {
			tb.AppendRow(
				table.Row{img.Id, img.Kind, img.Snippet.PlaylistId, img.Snippet.Type},
			)
		}
	}
	return err
}

func (pi *PlaylistImage) Insert(writer io.Writer) error {
	pi.EnsureService()
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

	switch pi.Output {
	case "json":
		utils.PrintJSON(res, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, pi.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "PlaylistImage inserted: %s\n", res.Id)
	}
	return nil
}

func (pi *PlaylistImage) Update(writer io.Writer) error {
	pi.EnsureService()
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

	switch pi.Output {
	case "json":
		utils.PrintJSON(res, pi.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, pi.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "PlaylistImage updated: %s\n", res.Id)
	}
	return nil
}

func (pi *PlaylistImage) Delete(writer io.Writer) error {
	pi.EnsureService()
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

func WithIds(ids []string) Option {
	return func(pi *PlaylistImage) {
		pi.Ids = ids
	}
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

func WithMaxResults(maxResults int64) Option {
	return func(pi *PlaylistImage) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		pi.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(pi *PlaylistImage) {
		pi.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(pi *PlaylistImage) {
		pi.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

var (
	WithParts    = common.WithParts[*PlaylistImage]
	WithOutput   = common.WithOutput[*PlaylistImage]
	WithJsonpath = common.WithJsonpath[*PlaylistImage]
	WithService  = common.WithService[*PlaylistImage]
)
