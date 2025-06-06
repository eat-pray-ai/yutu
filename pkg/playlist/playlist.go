package playlist

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"log"

	"google.golang.org/api/youtube/v3"
)

var (
	service           *youtube.Service
	errGetPlaylist    = errors.New("failed to get playlist")
	errInsertPlaylist = errors.New("failed to insert playlist")
	errUpdatePlaylist = errors.New("failed to update playlist")
	errDeletePlaylist = errors.New("failed to delete playlist")
)

type playlist struct {
	IDs         []string `yaml:"ids" json:"ids"`
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Hl          string   `yaml:"hl" json:"hl"`
	MaxResults  int64    `yaml:"max_results" json:"max_results"`
	Mine        *bool    `yaml:"mine" json:"mine"`
	Tags        []string `yaml:"tags" json:"tags"`
	Language    string   `yaml:"language" json:"language"`
	ChannelId   string   `yaml:"channel_id" json:"channel_id"`
	Privacy     string   `yaml:"privacy" json:"privacy"`

	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type Playlist interface {
	List([]string, string)
	Insert(output string)
	Update(output string)
	Delete()
	get([]string) []*youtube.Playlist
}

type Option func(*playlist)

func NewPlaylist(opts ...Option) Playlist {
	p := &playlist{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *playlist) get(parts []string) []*youtube.Playlist {
	call := service.Playlists.List(parts)

	if len(p.IDs) > 0 {
		call = call.Id(p.IDs...)
	}
	if p.Hl != "" {
		call = call.Hl(p.Hl)
	}
	if p.Mine != nil {
		call = call.Mine(*p.Mine)
	}
	if p.MaxResults <= 0 {
		p.MaxResults = 1
	}
	call = call.MaxResults(p.MaxResults)
	if p.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(p.OnBehalfOfContentOwner)
	}
	if p.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(p.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(p, nil)
		log.Fatalln(errors.Join(errGetPlaylist, err))
	}

	return res.Items
}

func (p *playlist) List(parts []string, output string) {
	playlists := p.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(playlists, nil)
	case "yaml":
		utils.PrintYAML(playlists, nil)
	default:
		fmt.Println("ID\tTitle")
		for _, playlist := range playlists {
			fmt.Printf("%s\t%s\n", playlist.Id, playlist.Snippet.Title)
		}
	}
}

func (p *playlist) Insert(output string) {
	upload := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:           p.Title,
			Description:     p.Description,
			Tags:            p.Tags,
			DefaultLanguage: p.Language,
			ChannelId:       p.ChannelId,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: p.Privacy,
		},
	}

	call := service.Playlists.Insert([]string{"snippet", "status"}, upload)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(p, nil)
		log.Fatalln(errors.Join(errInsertPlaylist, err))
	}

	switch output {
	case "json":
		utils.PrintJSON(res, nil)
	case "yaml":
		utils.PrintYAML(res, nil)
	case "silent":
	default:
		fmt.Printf("Playlist inserted: %s\n", res.Id)
	}
}

func (p *playlist) Update(output string) {
	playlist := p.get([]string{"id", "snippet", "status"})[0]
	if p.Title != "" {
		playlist.Snippet.Title = p.Title
	}
	if p.Description != "" {
		playlist.Snippet.Description = p.Description
	}
	if p.Tags != nil {
		playlist.Snippet.Tags = p.Tags
	}
	if p.Language != "" {
		playlist.Snippet.DefaultLanguage = p.Language
	}
	if p.Privacy != "" {
		playlist.Status.PrivacyStatus = p.Privacy
	}

	call := service.Playlists.Update([]string{"snippet", "status"}, playlist)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(p, nil)
		log.Fatalln(errors.Join(errUpdatePlaylist, err))
	}

	switch output {
	case "json":
		utils.PrintJSON(res, nil)
	case "yaml":
		utils.PrintYAML(res, nil)
	case "silent":
	default:
		fmt.Printf("Playlist updated: %s\n", res.Id)
	}
}

func (p *playlist) Delete() {
	for _, id := range p.IDs {
		call := service.Playlists.Delete(id)
		if p.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(p.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			utils.PrintJSON(p, nil)
			log.Fatalln(errors.Join(errDeletePlaylist, err))
		}
		fmt.Printf("Playlist %s deleted", id)
	}
}

func WithIDs(ids []string) Option {
	return func(p *playlist) {
		p.IDs = ids
	}
}

func WithTitle(title string) Option {
	return func(p *playlist) {
		p.Title = title
	}
}

func WithDescription(description string) Option {
	return func(p *playlist) {
		p.Description = description
	}
}

func WithTags(tags []string) Option {
	return func(p *playlist) {
		p.Tags = tags
	}
}

func WithLanguage(language string) Option {
	return func(p *playlist) {
		p.Language = language
	}
}

func WithChannelId(channelId string) Option {
	return func(p *playlist) {
		p.ChannelId = channelId
	}
}

func WithPrivacy(privacy string) Option {
	return func(p *playlist) {
		p.Privacy = privacy
	}
}

func WithHl(hl string) Option {
	return func(p *playlist) {
		p.Hl = hl
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(p *playlist) {
		p.MaxResults = maxResults
	}
}

func WithMine(mine *bool) Option {
	return func(p *playlist) {
		if mine != nil {
			p.Mine = mine
		}
	}
}

func WithOnBehalfOfContentOwner(contentOwner string) Option {
	return func(p *playlist) {
		p.OnBehalfOfContentOwner = contentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(channel string) Option {
	return func(p *playlist) {
		p.OnBehalfOfContentOwnerChannel = channel
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *playlist) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
