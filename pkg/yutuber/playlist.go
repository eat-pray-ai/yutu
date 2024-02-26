package yutuber

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylist    error = errors.New("failed to get playlist")
	errInsertPlaylist error = errors.New("failed to insert playlist")
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
	List()
	Insert()
	get(parts []string) []*youtube.Playlist
}

type PlaylistOption func(*Playlist)

func NewPlaylist(opts ...PlaylistOption) *Playlist {
	p := &Playlist{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Playlist) List() {
	playlists := p.get([]string{"snippet"})
	for _, playlist := range playlists {
		fmt.Printf("          ID: %s\n", playlist.Id)
		fmt.Printf("       Title: %s\n", playlist.Snippet.Title)
		fmt.Printf(" Description: %s\n", playlist.Snippet.Description)
		fmt.Printf("Published At: %s\n", playlist.Snippet.PublishedAt)
		fmt.Printf("     Channel: %s\n\n", playlist.Snippet.ChannelId)
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

func WithPlaylistTags(tags string) PlaylistOption {
	return func(p *Playlist) {
		p.tags = strings.Split(tags, ",")
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
