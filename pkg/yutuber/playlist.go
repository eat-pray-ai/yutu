package yutuber

import (
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
)

type Playlist struct {
	id           string
	title        string
	desc         string
	tags         string
	language     string
	thumbnail    string
	thumbnailVid string
	channelId    string
	privacy      string
}

type PlaylistService interface {
	List()
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
		log.Printf("          ID: %s\n", playlist.Id)
		log.Printf("       Title: %s\n", playlist.Snippet.Title)
		log.Printf(" Description: %s\n", playlist.Snippet.Description)
		log.Printf("Published At: %s\n", playlist.Snippet.PublishedAt)
		log.Printf("     Channel: %s\n\n", playlist.Snippet.ChannelId)
	}
}

func (p *Playlist) get(parts []string) []*youtube.Playlist {
	service := auth.NewY2BService(youtube.YoutubeReadonlyScope)
	call := service.Playlists.List(parts).Id(p.id)
	res, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call to get playlist: %v", err)
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
		p.tags = tags
	}
}

func WithPlaylistLanguage(language string) PlaylistOption {
	return func(p *Playlist) {
		p.language = language
	}
}

func WithPlaylistThumbnail(thumbnail string) PlaylistOption {
	return func(p *Playlist) {
		p.thumbnail = thumbnail
	}
}

func WithPlaylistThumbnailVid(thumbnailVid string) PlaylistOption {
	return func(p *Playlist) {
		p.thumbnailVid = thumbnailVid
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
