package yutuber

import (
	"errors"
	"fmt"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/utils"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylist    error = errors.New("failed to get playlist")
	errInsertPlaylist error = errors.New("failed to insert playlist")
	errUpdatePlaylist error = errors.New("failed to update playlist")
)

type Playlist struct {
	id        string
	title     string
	desc      string
	tags      []string
	language  string
	channelId string
	privacy   string
}

type PlaylistService interface {
	List([]string, string)
	Insert()
	Update()
	get([]string) []*youtube.Playlist
}

type PlaylistOption func(*Playlist)

func NewPlaylist(opts ...PlaylistOption) *Playlist {
	p := &Playlist{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Playlist) List(parts []string, output string) {
	playlists := p.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(playlists)
	case "yaml":
		utils.PrintYAML(playlists)
	default:
		fmt.Println("ID\tTitle")
		for _, playlist := range playlists {
			fmt.Printf("%s\t%s\n", playlist.Id, playlist.Snippet.Title)
		}
	}
}

func (p *Playlist) Insert() {
	upload := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:           p.title,
			Description:     p.desc,
			Tags:            p.tags,
			DefaultLanguage: p.language,
			ChannelId:       p.channelId,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: p.privacy,
		},
	}

	call := service.Playlists.Insert([]string{"snippet", "status"}, upload)
	playlist, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertPlaylist, err))
	}

	fmt.Printf("Playlist %s inserted\n", playlist.Id)
}

func (p *Playlist) Update() {
	playlist := p.get([]string{"id", "snippet", "status"})[0]
	if p.title != "" {
		playlist.Snippet.Title = p.title
	}
	if p.desc != "" {
		playlist.Snippet.Description = p.desc
	}
	if p.tags != nil {
		playlist.Snippet.Tags = p.tags
	}
	if p.language != "" {
		playlist.Snippet.DefaultLanguage = p.language
	}
	if p.privacy != "" {
		playlist.Status.PrivacyStatus = p.privacy
	}

	call := service.Playlists.Update([]string{"snippet", "status"}, playlist)
	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdatePlaylist, err), p.id)
	}

	data, _ := res.MarshalJSON()
	fmt.Println("Playlist updated:")
	utils.PrintJSON(data)
}

func (p *Playlist) get(parts []string) []*youtube.Playlist {
	call := service.Playlists.List(parts)
	switch {
	case p.id != "":
		call = call.Id(p.id)
	case p.channelId != "":
		call = call.ChannelId(p.channelId)
	default:
		call = call.Mine(true)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetPlaylist, err), p.id)
	}

	return res.Items
}

func WithPlaylistId(id string) PlaylistOption {
	return func(p *Playlist) {
		p.id = id
	}
}

func WithPlaylistTitle(title string) PlaylistOption {
	return func(p *Playlist) {
		p.title = title
	}
}

func WithPlaylistDesc(desc string) PlaylistOption {
	return func(p *Playlist) {
		p.desc = desc
	}
}

func WithPlaylistTags(tags []string) PlaylistOption {
	return func(p *Playlist) {
		p.tags = tags
	}
}

func WithPlaylistLanguage(language string) PlaylistOption {
	return func(p *Playlist) {
		p.language = language
	}
}

func WithPlaylistChannelId(channelId string) PlaylistOption {
	return func(p *Playlist) {
		p.channelId = channelId
	}
}

func WithPlaylistPrivacy(privacy string) PlaylistOption {
	return func(p *Playlist) {
		p.privacy = privacy
	}
}
