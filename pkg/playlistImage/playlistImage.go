package playlistImage

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
	"io"
	"os"
)

var (
	service                *youtube.Service
	errGetPlaylistImage    = errors.New("failed to get playlist image")
	errInsertPlaylistImage = errors.New("failed to insert playlist image")
	errUpdatePlaylistImage = errors.New("failed to update playlist image")
	errDeletePlaylistImage = errors.New("failed to delete playlist image")
)

type playlistImage struct {
	IDs        []string `yaml:"ids" json:"ids"`
	Height     int64    `yaml:"height" json:"height"`
	PlaylistID string   `yaml:"playlistId" json:"playlistId"`
	Type       string   `yaml:"type" json:"type"`
	Width      int64    `yaml:"width" json:"width"`
	File       string   `yaml:"file" json:"file"`

	Parent     string `yaml:"parent" json:"parent"`
	MaxResults int64  `yaml:"max_results" json:"max_results"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type PlaylistImage interface {
	Get([]string) ([]*youtube.PlaylistImage, error)
	List([]string, string, io.Writer) error
	Insert(string, io.Writer) error
	Update(string, io.Writer) error
	Delete(io.Writer) error
}

type Option func(*playlistImage)

func NewPlaylistImage(opts ...Option) PlaylistImage {
	pi := &playlistImage{}
	for _, opt := range opts {
		opt(pi)
	}
	return pi
}

func (pi *playlistImage) Get(parts []string) ([]*youtube.PlaylistImage, error) {
	call := service.PlaylistImages.List()
	call = call.Part(parts...)

	if pi.Parent != "" {
		call = call.Parent(pi.Parent)
	}
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(pi.OnBehalfOfContentOwnerChannel)
	}
	if pi.MaxResults <= 0 {
		pi.MaxResults = 1
	}
	call = call.MaxResults(pi.MaxResults)

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetPlaylistImage, err)
	}

	return res.Items, nil
}

func (pi *playlistImage) List(
	parts []string, output string, writer io.Writer,
) error {
	playlistImages, err := pi.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(playlistImages, writer)
	case "yaml":
		utils.PrintYAML(playlistImages, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Kind", "Playlist ID", "Type"})
		for _, img := range playlistImages {
			tb.AppendRow(
				table.Row{img.Id, img.Kind, img.Snippet.PlaylistId, img.Snippet.Type},
			)
		}
	}
	return nil
}

func (pi *playlistImage) Insert(output string, writer io.Writer) error {
	file, err := os.Open(pi.File)
	if err != nil {
		return errors.Join(errInsertPlaylistImage, err)
	}
	defer file.Close()

	playlistImage := &youtube.PlaylistImage{
		Kind: "youtube#playlistImages",
		Snippet: &youtube.PlaylistImageSnippet{
			PlaylistId: pi.PlaylistID,
			Type:       pi.Type,
			Height:     pi.Height,
			Width:      pi.Width,
		},
	}

	call := service.PlaylistImages.Insert(playlistImage)
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

	switch output {
	case "json":
		utils.PrintJSON(res, writer)
	case "yaml":
		utils.PrintYAML(res, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "PlaylistImage inserted: %s\n", res.Id)
	}
	return nil
}

func (pi *playlistImage) Update(output string, writer io.Writer) error {
	playlistImages, err := pi.Get([]string{"id", "kind", "snippet"})
	if err != nil {
		return errors.Join(errUpdatePlaylistImage, err)
	}
	if len(playlistImages) == 0 {
		return errGetPlaylistImage
	}

	playlistImage := playlistImages[0]
	if pi.PlaylistID != "" {
		playlistImage.Snippet.PlaylistId = pi.PlaylistID
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

	call := service.PlaylistImages.Update(playlistImage)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}
	if pi.File != "" {
		file, err := os.Open(pi.File)
		if err != nil {
			return errors.Join(errUpdatePlaylistImage, err)
		}
		defer file.Close()
		call = call.Media(file)
	}
	call = call.Part("id", "kind", "snippet")

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdatePlaylistImage, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, writer)
	case "yaml":
		utils.PrintYAML(res, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "PlaylistImage updated: %s\n", res.Id)
	}
	return nil
}

func (pi *playlistImage) Delete(writer io.Writer) error {
	for _, id := range pi.IDs {
		call := service.PlaylistImages.Delete()
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

func WithIDs(ids []string) Option {
	return func(pi *playlistImage) {
		pi.IDs = ids
	}
}

func WithHeight(height int64) Option {
	return func(pi *playlistImage) {
		pi.Height = height
	}
}

func WithPlaylistID(playlistID string) Option {
	return func(pi *playlistImage) {
		pi.PlaylistID = playlistID
	}
}

func WithType(t string) Option {
	return func(pi *playlistImage) {
		pi.Type = t
	}
}

func WithWidth(width int64) Option {
	return func(pi *playlistImage) {
		pi.Width = width
	}
}

func WithFile(file string) Option {
	return func(pi *playlistImage) {
		pi.File = file
	}
}

func WithParent(parent string) Option {
	return func(pi *playlistImage) {
		pi.Parent = parent
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(pi *playlistImage) {
		pi.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(pi *playlistImage) {
		pi.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(pi *playlistImage) {
		pi.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithService(svc *youtube.Service) Option {
	return func(pi *playlistImage) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
