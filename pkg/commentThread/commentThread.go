// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetCommentThread    = errors.New("failed to get comment thread")
	errInsertCommentThread = errors.New("failed to insert comment thread")
)

type CommentThread struct {
	*common.Fields
	Ids                          []string `yaml:"ids" json:"ids"`
	AllThreadsRelatedToChannelId string   `yaml:"all_threads_related_to_channel_id" json:"all_threads_related_to_channel_id"`
	AuthorChannelId              string   `yaml:"author_channel_id" json:"author_channel_id"`
	ChannelId                    string   `yaml:"channel_id" json:"channel_id"`
	MaxResults                   int64    `yaml:"max_results" json:"max_results"`
	ModerationStatus             string   `yaml:"moderation_status" json:"moderation_status"`
	Order                        string   `yaml:"order" json:"order"`
	SearchTerms                  string   `yaml:"search_terms" json:"search_terms"`
	TextFormat                   string   `yaml:"text_format" json:"text_format"`
	TextOriginal                 string   `yaml:"text_original" json:"text_original"`
	VideoId                      string   `yaml:"video_id" json:"video_id"`
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
	c.EnsureService()
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

	var items []*youtube.CommentThread
	pageToken := ""
	for c.MaxResults > 0 {
		call = call.MaxResults(min(c.MaxResults, pkg.PerPage))
		c.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetCommentThread, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (c *CommentThread) List(writer io.Writer) error {
	commentThreads, err := c.Get()
	if err != nil && commentThreads == nil {
		return err
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(commentThreads, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(commentThreads, c.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Author", "Video ID", "Text Display"})
		for _, cot := range commentThreads {
			snippet := cot.Snippet.TopLevelComment.Snippet
			tb.AppendRow(
				table.Row{
					cot.Id, snippet.AuthorDisplayName,
					snippet.VideoId, snippet.TextDisplay,
				},
			)
		}
	}
	return err
}

func (c *CommentThread) Insert(writer io.Writer) error {
	c.EnsureService()
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

	switch c.Output {
	case "json":
		utils.PrintJSON(res, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "CommentThread inserted: %s\n", res.Id)
	}
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

func WithChannelId(channelId string) Option {
	return func(c *CommentThread) {
		c.ChannelId = channelId
	}
}

func WithIds(ids []string) Option {
	return func(c *CommentThread) {
		c.Ids = ids
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *CommentThread) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		c.MaxResults = maxResults
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
	WithParts    = common.WithParts[*CommentThread]
	WithOutput   = common.WithOutput[*CommentThread]
	WithJsonpath = common.WithJsonpath[*CommentThread]
	WithService  = common.WithService[*CommentThread]
)
