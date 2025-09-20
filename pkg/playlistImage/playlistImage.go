package playlistImage

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
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

type PlaylistImage[T any] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Delete(io.Writer) error
}

type Option func(*playlistImage)

func NewPlaylistImage(opts ...Option) PlaylistImage[youtube.PlaylistImage] {
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

func (pi *playlistImage) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	playlistImages, err := pi.Get(parts)
	if err != nil && playlistImages == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(playlistImages, jpath, writer)
	case "yaml":
		utils.PrintYAML(playlistImages, jpath, writer)
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
	return err
}

func (pi *playlistImage) Insert(
	output string, jpath string, writer io.Writer,
) error {
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
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "PlaylistImage inserted: %s\n", res.Id)
	}
	return nil
}

func (pi *playlistImage) Update(
	output string, jpath string, writer io.Writer,
) error {
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

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
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
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
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
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
