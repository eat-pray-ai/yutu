// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

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
	service                *youtube.Service
	errGetCommentThread    = errors.New("failed to get comment thread")
	errInsertCommentThread = errors.New("failed to insert comment thread")
)

type commentThread struct {
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

type CommentThread[T any] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
	Insert(output string, s string, writer io.Writer) error
}

type Option func(*commentThread)

func NewCommentThread(opts ...Option) CommentThread[youtube.CommentThread] {
	c := &commentThread{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *commentThread) Get(parts []string) ([]*youtube.CommentThread, error) {
	call := service.CommentThreads.List(parts)
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

func (c *commentThread) List(
	parts []string, output string, jsonpath string, writer io.Writer,
) error {
	commentThreads, err := c.Get(parts)
	if err != nil && commentThreads == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(commentThreads, jsonpath, writer)
	case "yaml":
		utils.PrintYAML(commentThreads, jsonpath, writer)
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

func (c *commentThread) Insert(
	output string, jsonpath string, writer io.Writer,
) error {
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

	res, err := service.CommentThreads.Insert([]string{"snippet"}, ct).Do()
	if err != nil {
		return errors.Join(errInsertCommentThread, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "CommentThread inserted: %s\n", res.Id)
	}
	return nil
}

func WithAllThreadsRelatedToChannelId(allThreadsRelatedToChannelId string) Option {
	return func(c *commentThread) {
		c.AllThreadsRelatedToChannelId = allThreadsRelatedToChannelId
	}
}

func WithAuthorChannelId(authorChannelId string) Option {
	return func(c *commentThread) {
		c.AuthorChannelId = authorChannelId
	}
}

func WithChannelId(channelId string) Option {
	return func(c *commentThread) {
		c.ChannelId = channelId
	}
}

func WithIds(ids []string) Option {
	return func(c *commentThread) {
		c.Ids = ids
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *commentThread) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		c.MaxResults = maxResults
	}
}

func WithModerationStatus(moderationStatus string) Option {
	return func(c *commentThread) {
		c.ModerationStatus = moderationStatus
	}
}

func WithOrder(order string) Option {
	return func(c *commentThread) {
		c.Order = order
	}
}

func WithSearchTerms(searchTerms string) Option {
	return func(c *commentThread) {
		c.SearchTerms = searchTerms
	}
}

func WithTextFormat(textFormat string) Option {
	return func(c *commentThread) {
		c.TextFormat = textFormat
	}
}

func WithTextOriginal(textOriginal string) Option {
	return func(c *commentThread) {
		c.TextOriginal = textOriginal
	}
}

func WithVideoId(videoId string) Option {
	return func(c *commentThread) {
		c.VideoId = videoId
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *commentThread) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
