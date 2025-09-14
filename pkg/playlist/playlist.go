package playlist

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
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

type Playlist[T any] interface {
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Delete(io.Writer) error
	Get([]string) ([]*T, error)
}

type Option func(*playlist)

func NewPlaylist(opts ...Option) Playlist[youtube.Playlist] {
	p := &playlist{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *playlist) Get(parts []string) ([]*youtube.Playlist, error) {
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
	if p.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(p.OnBehalfOfContentOwner)
	}
	if p.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(p.OnBehalfOfContentOwnerChannel)
	}

	var items []*youtube.Playlist
	pageToken := ""
	for p.MaxResults > 0 {
		call = call.MaxResults(min(p.MaxResults, pkg.PerPage))
		p.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetPlaylist, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (p *playlist) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	playlists, err := p.Get(parts)
	if err != nil && playlists == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(playlists, jpath, writer)
	case "yaml":
		utils.PrintYAML(playlists, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Channel ID", "Title"})
		for _, pl := range playlists {
			tb.AppendRow(table.Row{pl.Id, pl.Snippet.ChannelId, pl.Snippet.Title})
		}
	}
	return err
}

func (p *playlist) Insert(output string, jpath string, writer io.Writer) error {
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
		return errors.Join(errInsertPlaylist, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist inserted: %s\n", res.Id)
	}
	return nil
}

func (p *playlist) Update(output string, jpath string, writer io.Writer) error {
	playlists, err := p.Get([]string{"id", "snippet", "status"})
	if err != nil {
		return errors.Join(errUpdatePlaylist, err)
	}
	if len(playlists) == 0 {
		return errGetPlaylist
	}

	playlist := playlists[0]
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
		return errors.Join(errUpdatePlaylist, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Playlist updated: %s\n", res.Id)
	}
	return nil
}

func (p *playlist) Delete(writer io.Writer) error {
	for _, id := range p.IDs {
		call := service.Playlists.Delete(id)
		if p.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(p.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeletePlaylist, err)
		}
		_, _ = fmt.Fprintf(writer, "Playlist %s deleted", id)
	}
	return nil
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
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
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
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
