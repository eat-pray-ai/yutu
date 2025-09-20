package caption

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service            *youtube.Service
	errGetCaption      = errors.New("failed to get caption")
	errUpdateCaption   = errors.New("failed to update caption")
	errDeleteCaption   = errors.New("failed to delete caption")
	errInsertCaption   = errors.New("failed to insert caption")
	errDownloadCaption = errors.New("failed to download caption")
)

type caption struct {
	IDs                    []string `yaml:"ids" json:"ids"`
	File                   string   `yaml:"file" json:"file"`
	AudioTrackType         string   `yaml:"audio_track_type" json:"audio_track_type"`
	IsAutoSynced           *bool    `yaml:"is_auto_synced" json:"is_auto_synced"`
	IsCC                   *bool    `yaml:"is_cc" json:"is_cc"`
	IsDraft                *bool    `yaml:"is_draft" json:"is_draft"`
	IsEasyReader           *bool    `yaml:"is_easy_reader" json:"is_easy_reader"`
	IsLarge                *bool    `yaml:"is_large" json:"is_large"`
	Language               string   `yaml:"language" json:"language"`
	Name                   string   `yaml:"name" json:"name"`
	TrackKind              string   `yaml:"track_kind" json:"track_kind"`
	OnBehalfOf             string   `yaml:"on_behalf_of" json:"on_behalf_of"`
	OnBehalfOfContentOwner string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	VideoId                string   `yaml:"video_id" json:"video_id"`
	Tfmt                   string   `yaml:"tfmt" json:"tfmt"`
	Tlang                  string   `yaml:"tlang" json:"tlang"`
}

type Caption[T youtube.Caption] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Delete(io.Writer) error
	Download(io.Writer) error
}

type Option func(*caption)

func NewCation(opts ...Option) Caption[youtube.Caption] {
	c := &caption{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *caption) Get(parts []string) ([]*youtube.Caption, error) {
	call := service.Captions.List(parts, c.VideoId)
	if len(c.IDs) > 0 {
		call = call.Id(c.IDs...)
	}
	if c.OnBehalfOf != "" {
		call = call.OnBehalfOf(c.OnBehalfOf)
	}
	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetCaption, err)
	}

	return res.Items, nil
}

func (c *caption) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	captions, err := c.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(captions, jpath, writer)
	case "yaml":
		utils.PrintYAML(captions, jpath, writer)
	case "table":
		tb := table.NewWriter()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Video ID", "Name", "Language"})
		defer tb.Render()
		for _, caption := range captions {
			tb.AppendRow(
				table.Row{
					caption.Id, caption.Snippet.VideoId,
					caption.Snippet.Name, caption.Snippet.Language,
				},
			)
		}
	}
	return nil
}

func (c *caption) Insert(
	output string, jpath string, writer io.Writer,
) error {
	file, err := pkg.Root.Open(c.File)
	if err != nil {
		return errors.Join(errInsertCaption, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	caption := &youtube.Caption{
		Snippet: &youtube.CaptionSnippet{
			AudioTrackType: c.AudioTrackType,
			IsAutoSynced:   *c.IsAutoSynced,
			IsCC:           *c.IsCC,
			IsDraft:        *c.IsDraft,
			IsEasyReader:   *c.IsEasyReader,
			IsLarge:        *c.IsLarge,
			Language:       c.Language,
			Name:           c.Name,
			TrackKind:      c.TrackKind,
			VideoId:        c.VideoId,
		},
	}

	call := service.Captions.Insert([]string{"snippet"}, caption).Media(file)
	if c.OnBehalfOf != "" {
		call = call.OnBehalfOf(c.OnBehalfOf)
	}
	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertCaption, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Caption inserted: %s\n", res.Id)
	}
	return nil
}

func (c *caption) Update(
	output string, jpath string, writer io.Writer,
) error {
	captions, err := c.Get([]string{"snippet"})
	if err != nil {
		return errors.Join(errUpdateCaption, err)
	}
	if len(captions) == 0 {
		return errGetCaption
	}

	caption := captions[0]
	if c.AudioTrackType != "" {
		caption.Snippet.AudioTrackType = c.AudioTrackType
	}
	if c.IsAutoSynced != nil {
		caption.Snippet.IsAutoSynced = *c.IsAutoSynced
	}
	if c.IsCC != nil {
		caption.Snippet.IsCC = *c.IsCC
	}
	if c.IsDraft != nil {
		caption.Snippet.IsDraft = *c.IsDraft
	}
	if c.IsEasyReader != nil {
		caption.Snippet.IsEasyReader = *c.IsEasyReader
	}
	if c.IsLarge != nil {
		caption.Snippet.IsLarge = *c.IsLarge
	}
	if c.Language != "" {
		caption.Snippet.Language = c.Language
	}
	if c.Name != "" {
		caption.Snippet.Name = c.Name
	}
	if c.TrackKind != "" {
		caption.Snippet.TrackKind = c.TrackKind
	}
	if c.VideoId != "" {
		caption.Snippet.VideoId = c.VideoId
	}

	call := service.Captions.Update([]string{"snippet"}, caption)
	if c.File != "" {
		file, err := pkg.Root.Open(c.File)
		if err != nil {
			return errors.Join(errUpdateCaption, err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		call = call.Media(file)
	}
	if c.OnBehalfOf != "" {
		call = call.OnBehalfOf(c.OnBehalfOf)
	}
	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateCaption, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Caption updated: %s\n", res.Id)
	}
	return nil
}

func (c *caption) Delete(writer io.Writer) error {
	for _, id := range c.IDs {
		call := service.Captions.Delete(id)
		if c.OnBehalfOf != "" {
			call = call.OnBehalfOf(c.OnBehalfOf)
		}
		if c.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteCaption, err)
		}

		_, _ = fmt.Fprintf(writer, "Caption %s deleted\n", id)
	}
	return nil
}

func (c *caption) Download(writer io.Writer) error {
	call := service.Captions.Download(c.IDs[0])
	if c.Tfmt != "" {
		call = call.Tfmt(c.Tfmt)
	}
	if c.Tlang != "" {
		call = call.Tlang(c.Tlang)
	}
	if c.OnBehalfOf != "" {
		call = call.OnBehalfOf(c.OnBehalfOf)
	}
	if c.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.OnBehalfOfContentOwner)
	}

	res, err := call.Download()
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}

	file, err := os.Create(c.File)
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = file.Write(body)
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}

	_, _ = fmt.Fprintf(writer, "Caption %s downloaded to %s\n", c.IDs[0], c.File)
	return nil
}

func WithIDs(ids []string) Option {
	return func(c *caption) {
		c.IDs = ids
	}
}

func WithFile(file string) Option {
	return func(c *caption) {
		c.File = file
	}
}

func WithAudioTrackType(audioTrackType string) Option {
	return func(c *caption) {
		c.AudioTrackType = audioTrackType
	}
}

func WithIsAutoSynced(isAutoSynced *bool) Option {
	return func(c *caption) {
		if isAutoSynced != nil {
			c.IsAutoSynced = isAutoSynced
		}
	}
}

func WithIsCC(isCC *bool) Option {
	return func(c *caption) {
		if isCC != nil {
			c.IsCC = isCC
		}
	}
}

func WithIsDraft(isDraft *bool) Option {
	return func(c *caption) {
		if isDraft != nil {
			c.IsDraft = isDraft
		}
	}
}

func WithIsEasyReader(isEasyReader *bool) Option {
	return func(c *caption) {
		if isEasyReader != nil {
			c.IsEasyReader = isEasyReader
		}
	}
}

func WithIsLarge(isLarge *bool) Option {
	return func(c *caption) {
		if isLarge != nil {
			c.IsLarge = isLarge
		}
	}
}

func WithLanguage(language string) Option {
	return func(c *caption) {
		c.Language = language
	}
}

func WithName(name string) Option {
	return func(c *caption) {
		c.Name = name
	}
}

func WithTrackKind(trackKind string) Option {
	return func(c *caption) {
		c.TrackKind = trackKind
	}
}

func WithOnBehalfOf(onBehalfOf string) Option {
	return func(c *caption) {
		c.OnBehalfOf = onBehalfOf
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(c *caption) {
		c.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithVideoId(videoId string) Option {
	return func(c *caption) {
		c.VideoId = videoId
	}
}

func WithTfmt(tfmt string) Option {
	return func(c *caption) {
		c.Tfmt = tfmt
	}
}

func WithTlang(tlang string) Option {
	return func(c *caption) {
		c.Tlang = tlang
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *caption) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
