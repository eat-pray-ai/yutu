package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"log"

	"google.golang.org/api/youtube/v3"
)

var (
	errGetPlaylist    = errors.New("failed to get playlist")
	errInsertPlaylist = errors.New("failed to insert playlist")
	errUpdatePlaylist = errors.New("failed to update playlist")
)

type playlist struct {
	id         string
	title      string
	desc       string
	hl         string
	maxResults int64
	mine       string
	tags       []string
	language   string
	channelId  string
	privacy    string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
}

type Playlist interface {
	List([]string, string)
	Insert()
	Update()
	get([]string) []*youtube.Playlist
}

type PlaylistOption func(*playlist)

func NewPlaylist(opts ...PlaylistOption) Playlist {
	p := &playlist{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *playlist) get(parts []string) []*youtube.Playlist {
	call := service.Playlists.List(parts)

	if p.id != "" {
		call = call.Id(p.id)
	}
	if p.hl != "" {
		call = call.Hl(p.hl)
	}
	if p.mine == "true" {
		call = call.Mine(true)
	} else if p.mine == "false" {
		call = call.Mine(false)
	}

	call = call.MaxResults(p.maxResults)
	if p.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(p.onBehalfOfContentOwner)
	}
	if p.onBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(p.onBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetPlaylist, err), p.id)
	}

	return res.Items
}

func (p *playlist) List(parts []string, output string) {
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

func (p *playlist) Insert() {
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

func (p *playlist) Update() {
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

func WithPlaylistId(id string) PlaylistOption {
	return func(p *playlist) {
		p.id = id
	}
}

func WithPlaylistTitle(title string) PlaylistOption {
	return func(p *playlist) {
		p.title = title
	}
}

func WithPlaylistDesc(desc string) PlaylistOption {
	return func(p *playlist) {
		p.desc = desc
	}
}

func WithPlaylistTags(tags []string) PlaylistOption {
	return func(p *playlist) {
		p.tags = tags
	}
}

func WithPlaylistLanguage(language string) PlaylistOption {
	return func(p *playlist) {
		p.language = language
	}
}

func WithPlaylistChannelId(channelId string) PlaylistOption {
	return func(p *playlist) {
		p.channelId = channelId
	}
}

func WithPlaylistPrivacy(privacy string) PlaylistOption {
	return func(p *playlist) {
		p.privacy = privacy
	}
}

func WithPlaylistHl(hl string) PlaylistOption {
	return func(p *playlist) {
		p.hl = hl
	}
}

func WithPlaylistMaxResults(maxResults int64) PlaylistOption {
	return func(p *playlist) {
		p.maxResults = maxResults
	}
}

func WithPlaylistMine(mine string) PlaylistOption {
	return func(p *playlist) {
		p.mine = mine
	}
}

func WithPlaylistOnBehalfOfContentOwner(contentOwner string) PlaylistOption {
	return func(p *playlist) {
		p.onBehalfOfContentOwner = contentOwner
	}
}

func WithPlaylistOnBehalfOfContentOwnerChannel(channel string) PlaylistOption {
	return func(p *playlist) {
		p.onBehalfOfContentOwnerChannel = channel
	}
}
