package playlistImage

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
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
	ID         string `yaml:"id" json:"id"`
	Kind       string `yaml:"kind" json:"kind"`
	Height     int64  `yaml:"height" json:"height"`
	PlaylistID string `yaml:"playlistId" json:"playlistId"`
	Type       string `yaml:"type" json:"type"`
	Width      int64  `yaml:"width" json:"width"`
	File       string `yaml:"file" json:"file"`

	Parent     string `yaml:"parent" json:"parent"`
	MaxResults int64  `yaml:"max_results" json:"max_results"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type PlaylistImage interface {
	get([]string) []*youtube.PlaylistImage
	List([]string, string)
	Insert(string)
	Update(string)
	Delete()
}

type Option func(*playlistImage)

func NewPlaylistImage(opts ...Option) PlaylistImage {
	pi := &playlistImage{}
	for _, opt := range opts {
		opt(pi)
	}
	return pi
}

func (pi *playlistImage) get(parts []string) []*youtube.PlaylistImage {
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
		utils.PrintJSON(pi, nil)
		log.Fatalln(errors.Join(errGetPlaylistImage, err))
	}

	return res.Items
}

func (pi *playlistImage) List(parts []string, output string) {
	playlistImages := pi.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(playlistImages, nil)
	case "yaml":
		utils.PrintYAML(playlistImages, nil)
	default:
		fmt.Println("ID\tKind\tPlaylistID\tType")
		for _, image := range playlistImages {
			fmt.Printf(
				"%s\t%s\t%s\t%s\n",
				image.Id, image.Kind, image.Snippet.PlaylistId, image.Snippet.Type,
			)
		}
	}
}

func (pi *playlistImage) Insert(output string) {
	file, err := os.Open(pi.File)
	if err != nil {
		utils.PrintJSON(pi, nil)
		log.Fatalln(errors.Join(errInsertPlaylistImage, err))
	}
	defer file.Close()

	playlistImage := &youtube.PlaylistImage{
		Kind: pi.Kind,
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
		utils.PrintJSON(pi, nil)
		log.Fatalln(errors.Join(errInsertPlaylistImage, err))
	}

	switch output {
	case "json":
		utils.PrintJSON(res, nil)
	case "yaml":
		utils.PrintYAML(res, nil)
	default:
		fmt.Printf("PlaylistImage inserted: %s\n", res.Id)
	}
}

func (pi *playlistImage) Update(output string) {
	playlistImage := pi.get([]string{"id", "kind", "snippet"})[0]
	if pi.Kind != "" {
		playlistImage.Kind = pi.Kind
	}
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
			utils.PrintJSON(pi, nil)
			log.Fatalln(errors.Join(errUpdatePlaylistImage, err))
		}
		defer file.Close()
		call = call.Media(file)
	}
	call = call.Part("id", "kind", "snippet")

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(pi, nil)
		log.Fatalln(errors.Join(errUpdatePlaylistImage, err))
	}

	switch output {
	case "json":
		utils.PrintJSON(res, nil)
	case "yaml":
		utils.PrintYAML(res, nil)
	default:
		fmt.Printf("PlaylistImage updated: %s\n", res.Id)
	}
}

func (pi *playlistImage) Delete() {
	call := service.PlaylistImages.Delete()
	call = call.Id(pi.ID)
	if pi.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(pi.OnBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		utils.PrintJSON(pi, nil)
		log.Fatalln(errors.Join(errDeletePlaylistImage, err))
	}
	fmt.Printf("PlaylistImage %s deleted\n", pi.ID)
}

func WithID(id string) Option {
	return func(pi *playlistImage) {
		pi.ID = id
	}
}

func WithKind(kind string) Option {
	return func(pi *playlistImage) {
		pi.Kind = kind
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
