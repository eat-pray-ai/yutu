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
	id          string
	title       string
	description string
	hl          string
	maxResults  int64
	mine        string
	tags        []string
	language    string
	channelId   string
	privacy     string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
}

type Playlist interface {
	List([]string, string)
	Insert()
	Update()
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
			Description:     p.description,
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
	if p.description != "" {
		playlist.Snippet.Description = p.description
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

	fmt.Println("Playlist updated:")
	utils.PrintYAML(res)
}

func (p *playlist) Delete() {
	call := service.Playlists.Delete(p.id)
	if p.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(p.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeletePlaylist, err), p.id)
	}
	fmt.Printf("Playlist %s deleted", p.id)
}

func WithId(id string) Option {
	return func(p *playlist) {
		p.id = id
	}
}

func WithTitle(title string) Option {
	return func(p *playlist) {
		p.title = title
	}
}

func WithDescription(description string) Option {
	return func(p *playlist) {
		p.description = description
	}
}

func WithTags(tags []string) Option {
	return func(p *playlist) {
		p.tags = tags
	}
}

func WithLanguage(language string) Option {
	return func(p *playlist) {
		p.language = language
	}
}

func WithChannelId(channelId string) Option {
	return func(p *playlist) {
		p.channelId = channelId
	}
}

func WithPrivacy(privacy string) Option {
	return func(p *playlist) {
		p.privacy = privacy
	}
}

func WithHl(hl string) Option {
	return func(p *playlist) {
		p.hl = hl
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(p *playlist) {
		p.maxResults = maxResults
	}
}

func WithMine(mine string) Option {
	return func(p *playlist) {
		p.mine = mine
	}
}

func WithOnBehalfOfContentOwner(contentOwner string) Option {
	return func(p *playlist) {
		p.onBehalfOfContentOwner = contentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(channel string) Option {
	return func(p *playlist) {
		p.onBehalfOfContentOwnerChannel = channel
	}
}

func WithService() Option {
	return func(p *playlist) {
		service = auth.NewY2BService()
	}
}
