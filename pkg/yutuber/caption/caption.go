package caption

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service          *youtube.Service
	errOpenFile      = errors.New("failed to open file")
	errGetCaption    = errors.New("failed to get caption")
	errUpdateCaption = errors.New("failed to update caption")
	errDeleteCaption = errors.New("failed to delete caption")
	errInsertCaption = errors.New("failed to insert caption")
)

type caption struct {
	id                     string
	file                   string
	audioTrackType         string
	isAutoSynced           *bool
	isCC                   *bool
	isDraft                *bool
	isEasyReader           *bool
	isLarge                *bool
	language               string
	name                   string
	trackKind              string
	onBehalfOf             string
	onBehalfOfContentOwner string
	videoId                string
}

type Caption interface {
	get(parts []string) []*youtube.Caption
	List(parts []string, output string)
	Insert()
	Update()
	Delete()
}

type Option func(*caption)

func NewCation(opts ...Option) Caption {
	c := &caption{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *caption) get(parts []string) []*youtube.Caption {
	call := service.Captions.List(parts, c.videoId)
	if c.id != "" {
		call = call.Id(c.id)
	}
	if c.onBehalfOf != "" {
		call = call.OnBehalfOf(c.onBehalfOf)
	}
	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetCaption, err))
	}

	return res.Items
}

func (c *caption) List(parts []string, output string) {
	captions := c.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(captions)
	case "yaml":
		utils.PrintYAML(captions)
	default:
		fmt.Println("ID\tName")
		for _, caption := range captions {
			fmt.Printf("%s\t%s\n", caption.Id, caption.Snippet.Name)
		}
	}
}

func (c *caption) Insert() {
	file, err := os.Open(c.file)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err))
	}
	defer file.Close()

	caption := &youtube.Caption{
		Snippet: &youtube.CaptionSnippet{
			AudioTrackType: c.audioTrackType,
			IsAutoSynced:   *c.isAutoSynced,
			IsCC:           *c.isCC,
			IsDraft:        *c.isDraft,
			IsEasyReader:   *c.isEasyReader,
			IsLarge:        *c.isLarge,
			Language:       c.language,
			Name:           c.name,
			TrackKind:      c.trackKind,
			VideoId:        c.videoId,
		},
	}

	call := service.Captions.Insert([]string{"snippet"}, caption).Media(file)
	if c.onBehalfOf != "" {
		call = call.OnBehalfOf(c.onBehalfOf)
	}
	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errInsertCaption, err))
	}
	fmt.Printf("Caption %s inserted\n", res.Id)
}

func (c *caption) Update() {
	caption := c.get([]string{"snippet"})[0]
	if c.audioTrackType != "" {
		caption.Snippet.AudioTrackType = c.audioTrackType
	}
	if c.isAutoSynced != nil {
		caption.Snippet.IsAutoSynced = *c.isAutoSynced
	}
	if c.isCC != nil {
		caption.Snippet.IsCC = *c.isCC
	}
	if c.isDraft != nil {
		caption.Snippet.IsDraft = *c.isDraft
	}
	if c.isEasyReader != nil {
		caption.Snippet.IsEasyReader = *c.isEasyReader
	}
	if c.isLarge != nil {
		caption.Snippet.IsLarge = *c.isLarge
	}
	if c.language != "" {
		caption.Snippet.Language = c.language
	}
	if c.name != "" {
		caption.Snippet.Name = c.name
	}
	if c.trackKind != "" {
		caption.Snippet.TrackKind = c.trackKind
	}
	if c.videoId != "" {
		caption.Snippet.VideoId = c.videoId
	}

	call := service.Captions.Update([]string{"snippet"}, caption)
	if c.file != "" {
		file, err := os.Open(c.file)
		if err != nil {
			log.Fatalln(errors.Join(errOpenFile, err))
		}
		defer file.Close()
		call = call.Media(file)
	}
	if c.onBehalfOf != "" {
		call = call.OnBehalfOf(c.onBehalfOf)
	}
	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUpdateCaption, err))
	}

	fmt.Printf("Caption %s updated\n", res.Id)
}

func (c *caption) Delete() {
	call := service.Captions.Delete(c.id)
	if c.onBehalfOf != "" {
		call = call.OnBehalfOf(c.onBehalfOf)
	}
	if c.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(c.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errDeleteCaption, err))
	}

	fmt.Printf("Caption %s deleted\n", c.id)
}

func WithId(id string) Option {
	return func(c *caption) {
		c.id = id
	}
}

func WithFile(file string) Option {
	return func(c *caption) {
		c.file = file
	}
}

func WithAudioTrackType(audioTrackType string) Option {
	return func(c *caption) {
		c.audioTrackType = audioTrackType
	}
}

func WithIsAutoSynced(isAutoSynced bool, changed bool) Option {
	return func(c *caption) {
		if changed {
			c.isAutoSynced = &isAutoSynced
		}
	}
}

func WithIsCC(isCC bool, changed bool) Option {
	return func(c *caption) {
		if changed {
			c.isCC = &isCC
		}
	}
}

func WithIsDraft(isDraft bool, changed bool) Option {
	return func(c *caption) {
		if changed {
			c.isDraft = &isDraft
		}
	}
}

func WithIsEasyReader(isEasyReader bool, changed bool) Option {
	return func(c *caption) {
		if changed {
			c.isEasyReader = &isEasyReader
		}
	}
}

func WithIsLarge(isLarge bool, changed bool) Option {
	return func(c *caption) {
		if changed {
			c.isLarge = &isLarge
		}
	}
}

func WithLanguage(language string) Option {
	return func(c *caption) {
		c.language = language
	}
}

func WithName(name string) Option {
	return func(c *caption) {
		c.name = name
	}
}

func WithTrackKind(trackKind string) Option {
	return func(c *caption) {
		c.trackKind = trackKind
	}
}

func WithOnBehalfOf(onBehalfOf string) Option {
	return func(c *caption) {
		c.onBehalfOf = onBehalfOf
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(c *caption) {
		c.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithVideoId(videoId string) Option {
	return func(c *caption) {
		c.videoId = videoId
	}
}
