// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetCommentThread    = errors.New("failed to get comment thread")
	errInsertCommentThread = errors.New("failed to insert comment thread")
)

type CommentThread struct {
	*common.Fields
	AuthorChannelId  string `yaml:"author_channel_id" json:"author_channel_id,omitempty"`
	ModerationStatus string `yaml:"moderation_status" json:"moderation_status,omitempty"`
	Order            string `yaml:"order" json:"order,omitempty"`
	SearchTerms      string `yaml:"search_terms" json:"search_terms,omitempty"`
	TextFormat       string `yaml:"text_format" json:"text_format,omitempty"`
	TextOriginal     string `yaml:"text_original" json:"text_original,omitempty"`
	VideoId          string `yaml:"video_id" json:"video_id,omitempty"`

	AllThreadsRelatedToChannelId string `yaml:"all_threads_related_to_channel_id" json:"all_threads_related_to_channel_id,omitempty"`
}

type ICommentThread[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Insert(io.Writer) error
}

type Option func(*CommentThread)

func NewCommentThread(opts ...Option) ICommentThread[youtube.CommentThread] {
	c := &CommentThread{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *CommentThread) Get() ([]*youtube.CommentThread, error) {
	if err := c.EnsureService(); err != nil {
		return nil, err
	}
	call := c.Service.CommentThreads.List(c.Parts)
	if len(c.Ids) > 0 {
		call = call.Id(c.Ids...)
	}
	if c.AllThreadsRelatedToChannelId != "" {
		call = call.AllThreadsRelatedToChannelId(c.AllThreadsRelatedToChannelId)
	}
	if c.ChannelId != "" {
		call = call.ChannelId(c.ChannelId)
	}
	if c.ModerationStatus != "" {
		call = call.ModerationStatus(c.ModerationStatus)
	}
	if c.Order != "" {
		call = call.Order(c.Order)
	}
	if c.SearchTerms != "" {
		call = call.SearchTerms(c.SearchTerms)
	}
	if c.TextFormat != "" {
		call = call.TextFormat(c.TextFormat)
	}
	if c.VideoId != "" {
		call = call.VideoId(c.VideoId)
	}

	return common.Paginate(c.Fields, call, func(r *youtube.CommentThreadListResponse) ([]*youtube.CommentThread, string) {
		return r.Items, r.NextPageToken
	}, errGetCommentThread)
}

func (c *CommentThread) List(writer io.Writer) error {
	commentThreads, err := c.Get()
	if err != nil && commentThreads == nil {
		return err
	}

	common.PrintList(c.Output, commentThreads, writer, table.Row{"ID", "Author", "Video ID", "Text Display"}, func(cot *youtube.CommentThread) table.Row {
		snippet := cot.Snippet.TopLevelComment.Snippet
		return table.Row{cot.Id, snippet.AuthorDisplayName, snippet.VideoId, snippet.TextDisplay}
	})
	return err
}

func (c *CommentThread) Insert(writer io.Writer) error {
	if err := c.EnsureService(); err != nil {
		return err
	}
	ct := &youtube.CommentThread{
		Snippet: &youtube.CommentThreadSnippet{
			ChannelId: c.ChannelId,
			TopLevelComment: &youtube.Comment{
				Snippet: &youtube.CommentSnippet{
					AuthorChannelId: &youtube.CommentSnippetAuthorChannelId{
						Value: c.AuthorChannelId,
					},
					ChannelId:    c.ChannelId,
					TextOriginal: c.TextOriginal,
					VideoId:      c.VideoId,
				},
			},
		},
	}

	res, err := c.Service.CommentThreads.Insert([]string{"snippet"}, ct).Do()
	if err != nil {
		return errors.Join(errInsertCommentThread, err)
	}

	common.PrintResult(c.Output, res, writer, "CommentThread inserted: %s\n", res.Id)
	return nil
}

func WithAllThreadsRelatedToChannelId(allThreadsRelatedToChannelId string) Option {
	return func(c *CommentThread) {
		c.AllThreadsRelatedToChannelId = allThreadsRelatedToChannelId
	}
}

func WithAuthorChannelId(authorChannelId string) Option {
	return func(c *CommentThread) {
		c.AuthorChannelId = authorChannelId
	}
}

func WithModerationStatus(moderationStatus string) Option {
	return func(c *CommentThread) {
		c.ModerationStatus = moderationStatus
	}
}

func WithOrder(order string) Option {
	return func(c *CommentThread) {
		c.Order = order
	}
}

func WithSearchTerms(searchTerms string) Option {
	return func(c *CommentThread) {
		c.SearchTerms = searchTerms
	}
}

func WithTextFormat(textFormat string) Option {
	return func(c *CommentThread) {
		c.TextFormat = textFormat
	}
}

func WithTextOriginal(textOriginal string) Option {
	return func(c *CommentThread) {
		c.TextOriginal = textOriginal
	}
}

func WithVideoId(videoId string) Option {
	return func(c *CommentThread) {
		c.VideoId = videoId
	}
}

var (
	WithParts      = common.WithParts[*CommentThread]
	WithOutput     = common.WithOutput[*CommentThread]
	WithService    = common.WithService[*CommentThread]
	WithIds        = common.WithIds[*CommentThread]
	WithMaxResults = common.WithMaxResults[*CommentThread]
	WithChannelId  = common.WithChannelId[*CommentThread]
)
