package comment

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service                *youtube.Service
	errGetComment          = errors.New("failed to get comment")
	errMarkAsSpam          = errors.New("failed to mark comment as spam")
	errDeleteComment       = errors.New("failed to delete comment")
	errInsertComment       = errors.New("failed to insert comment")
	errUpdateComment       = errors.New("failed to update comment")
	errSetModerationStatus = errors.New("failed to set comment moderation status")
)

type comment struct {
	IDs              []string `yaml:"ids" json:"ids"`
	AuthorChannelId  string   `yaml:"author_channel_id" json:"author_channel_id"`
	CanRate          *bool    `yaml:"can_rate" json:"can_rate"`
	ChannelId        string   `yaml:"channel_id" json:"channel_id"`
	MaxResults       int64    `yaml:"max_results" json:"max_results"`
	ParentId         string   `yaml:"parent_id" json:"parent_id"`
	TextFormat       string   `yaml:"text_format" json:"text_format"`
	TextOriginal     string   `yaml:"text_original" json:"text_original"`
	ModerationStatus string   `yaml:"moderation_status" json:"moderation_status"`
	BanAuthor        *bool    `yaml:"ban_author" json:"ban_author"`
	VideoId          string   `yaml:"video_id" json:"video_id"`
	ViewerRating     string   `yaml:"viewer_rating" json:"viewer_rating"`
}

type Comment interface {
	get([]string) []*youtube.Comment
	List([]string, string)
	Insert(silent bool)
	Update(silent bool)
	Delete()
	MarkAsSpam(silent bool)
	SetModerationStatus(silent bool)
}

type Option func(*comment)

func NewComment(opts ...Option) Comment {
	c := &comment{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *comment) get(parts []string) []*youtube.Comment {
	call := service.Comments.List(parts)
	var result []*youtube.Comment

	for _, id := range c.IDs {
		if c.IDs[0] != "" {
			call = call.Id(id)
		}

		if c.MaxResults <= 0 {
			c.MaxResults = 1
		}
		call = call.MaxResults(c.MaxResults)

		if c.ParentId != "" {
			call = call.ParentId(c.ParentId)
		}

		if c.TextFormat != "" {
			call = call.TextFormat(c.TextFormat)
		}

		res, err := call.Do()
		if err != nil {
			utils.PrintJSON(c)
			log.Fatalln(errors.Join(errGetComment, err))
		}

		result = append(result, res.Items...)
	}

	return result
}

func (c *comment) List(parts []string, output string) {
	comments := c.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(comments)
	case "yaml":
		utils.PrintYAML(comments)
	default:
		fmt.Println("ID\tTextDisplay")
		for _, comment := range comments {
			fmt.Printf("%s\t%s\n", comment.Id, comment.Snippet.TextDisplay)
		}
	}
}

func (c *comment) Insert(silent bool) {
	comment := &youtube.Comment{
		Snippet: &youtube.CommentSnippet{
			AuthorChannelId: &youtube.CommentSnippetAuthorChannelId{
				Value: c.AuthorChannelId,
			},
			ChannelId:    c.ChannelId,
			ParentId:     c.ParentId,
			TextOriginal: c.TextOriginal,
			VideoId:      c.VideoId,
		},
	}

	call := service.Comments.Insert([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errInsertComment, err))
	}

	if !silent {
		fmt.Printf("Comment %s inserted", res.Id)
	}
}

func (c *comment) Update(silent bool) {
	comment := c.get([]string{"id", "snippet"})[0]

	if c.CanRate != nil {
		comment.Snippet.CanRate = *c.CanRate
	}

	if c.TextOriginal != "" {
		comment.Snippet.TextOriginal = c.TextOriginal
	}

	if c.ViewerRating != "" {
		comment.Snippet.ViewerRating = c.ViewerRating
	}

	call := service.Comments.Update([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errUpdateComment, err))
	}

	if !silent {
		fmt.Printf("Comment %s updated", res.Id)
	}
}

func (c *comment) MarkAsSpam(silent bool) {
	call := service.Comments.MarkAsSpam(c.IDs)
	err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errMarkAsSpam, err))
	}

	if !silent {
		fmt.Printf("Comment %s marked as spam", c.IDs)
	}
}

func (c *comment) SetModerationStatus(silent bool) {
	call := service.Comments.SetModerationStatus(c.IDs, c.ModerationStatus)

	if c.BanAuthor != nil {
		call = call.BanAuthor(*c.BanAuthor)
	}

	err := call.Do()
	if err != nil {
		utils.PrintJSON(c)
		log.Fatalln(errors.Join(errSetModerationStatus, err))
	}

	if !silent {
		fmt.Printf("Comment %s moderation status set to %s", c.IDs, c.ModerationStatus)
	}
}

func (c *comment) Delete() {
	for _, id := range c.IDs {
		call := service.Comments.Delete(id)
		err := call.Do()
		if err != nil {
			utils.PrintJSON(c)
			log.Fatalln(errors.Join(errDeleteComment, err))
		}
	}
}

func WithIDs(ids []string) Option {
	return func(c *comment) {
		c.IDs = ids
	}
}

func WithAuthorChannelId(authorChannelId string) Option {
	return func(c *comment) {
		c.AuthorChannelId = authorChannelId
	}
}

func WithCanRate(canRate bool, changed bool) Option {
	return func(c *comment) {
		if changed {
			c.CanRate = &canRate
		}
	}
}

func WithChannelId(channelId string) Option {
	return func(c *comment) {
		c.ChannelId = channelId
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *comment) {
		c.MaxResults = maxResults
	}
}

func WithParentId(parentId string) Option {
	return func(c *comment) {
		c.ParentId = parentId
	}
}

func WithTextFormat(textFormat string) Option {
	return func(c *comment) {
		c.TextFormat = textFormat
	}
}

func WithTextOriginal(textOriginal string) Option {
	return func(c *comment) {
		c.TextOriginal = textOriginal
	}
}

func WithModerationStatus(moderationStatus string) Option {
	return func(c *comment) {
		c.ModerationStatus = moderationStatus
	}
}

func WithBanAuthor(banAuthor bool, changed bool) Option {
	return func(c *comment) {
		if changed {
			c.BanAuthor = &banAuthor
		}
	}
}

func WithVideoId(videoId string) Option {
	return func(c *comment) {
		c.VideoId = videoId
	}
}

func WithViewerRating(viewerRating string) Option {
	return func(c *comment) {
		c.ViewerRating = viewerRating
	}
}
