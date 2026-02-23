// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

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
	errGetComment          = errors.New("failed to get comment")
	errMarkAsSpam          = errors.New("failed to mark comment as spam")
	errDeleteComment       = errors.New("failed to delete comment")
	errInsertComment       = errors.New("failed to insert comment")
	errUpdateComment       = errors.New("failed to update comment")
	errSetModerationStatus = errors.New("failed to set comment moderation status")
)

type Comment struct {
	*common.Fields
	Ids              []string `yaml:"ids" json:"ids,omitempty"`
	AuthorChannelId  string   `yaml:"author_channel_id" json:"author_channel_id,omitempty"`
	CanRate          *bool    `yaml:"can_rate" json:"can_rate,omitempty"`
	ChannelId        string   `yaml:"channel_id" json:"channel_id,omitempty"`
	MaxResults       int64    `yaml:"max_results" json:"max_results,omitempty"`
	ParentId         string   `yaml:"parent_id" json:"parent_id,omitempty"`
	TextFormat       string   `yaml:"text_format" json:"text_format,omitempty"`
	TextOriginal     string   `yaml:"text_original" json:"text_original,omitempty"`
	ModerationStatus string   `yaml:"moderation_status" json:"moderation_status,omitempty"`
	BanAuthor        *bool    `yaml:"ban_author" json:"ban_author,omitempty"`
	VideoId          string   `yaml:"video_id" json:"video_id,omitempty"`
	ViewerRating     string   `yaml:"viewer_rating" json:"viewer_rating,omitempty"`
}

type IComment[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
	MarkAsSpam(io.Writer) error
	SetModerationStatus(io.Writer) error
}

type Option func(*Comment)

func NewComment(opts ...Option) IComment[youtube.Comment] {
	c := &Comment{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Comment) Get() ([]*youtube.Comment, error) {
	c.EnsureService()
	call := c.Service.Comments.List(c.Parts)
	if len(c.Ids) > 0 && c.Ids[0] != "" {
		call = call.Id(c.Ids...)
	}
	if c.ParentId != "" {
		call = call.ParentId(c.ParentId)
	}
	if c.TextFormat != "" {
		call = call.TextFormat(c.TextFormat)
	}

	var items []*youtube.Comment
	pageToken := ""
	for c.MaxResults > 0 {
		call = call.MaxResults(min(c.MaxResults, pkg.PerPage))
		c.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetComment, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (c *Comment) List(writer io.Writer) error {
	comments, err := c.Get()
	if err != nil && comments == nil {
		return err
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(comments, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(comments, c.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Author", "Video ID", "Text Display"})
		for _, comment := range comments {
			tb.AppendRow(
				table.Row{
					comment.Id, comment.Snippet.AuthorDisplayName,
					comment.Snippet.VideoId, comment.Snippet.TextDisplay,
				},
			)
		}
	}
	return err
}

func (c *Comment) Insert(writer io.Writer) error {
	c.EnsureService()
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

	if c.CanRate != nil {
		comment.Snippet.CanRate = *c.CanRate
	}

	call := c.Service.Comments.Insert([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertComment, err)
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(res, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment inserted: %s\n", res.Id)
	}
	return nil
}

func (c *Comment) Update(writer io.Writer) error {
	c.EnsureService()
	c.Parts = []string{"id", "snippet"}
	comments, err := c.Get()

	if err != nil {
		return errors.Join(errUpdateComment, err)
	}
	if len(comments) == 0 {
		return errGetComment
	}

	comment := comments[0]
	if c.CanRate != nil {
		comment.Snippet.CanRate = *c.CanRate
	}

	if c.TextOriginal != "" {
		comment.Snippet.TextOriginal = c.TextOriginal
	}

	if c.ViewerRating != "" {
		comment.Snippet.ViewerRating = c.ViewerRating
	}

	call := c.Service.Comments.Update([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateComment, err)
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(res, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment updated: %s\n", res.Id)
	}
	return nil
}

func (c *Comment) MarkAsSpam(writer io.Writer) error {
	c.EnsureService()
	call := c.Service.Comments.MarkAsSpam(c.Ids)
	err := call.Do()
	if err != nil {
		return errors.Join(errMarkAsSpam, err)
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(c, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(c, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment marked as spam: %s\n", c.Ids)
	}
	return nil
}

func (c *Comment) SetModerationStatus(writer io.Writer) error {
	c.EnsureService()
	call := c.Service.Comments.SetModerationStatus(c.Ids, c.ModerationStatus)

	if c.BanAuthor != nil {
		call = call.BanAuthor(*c.BanAuthor)
	}

	err := call.Do()
	if err != nil {
		return errors.Join(errSetModerationStatus, err)
	}

	switch c.Output {
	case "json":
		utils.PrintJSON(c, c.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(c, c.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(
			writer, "Comment moderation status set to %s: %s\n",
			c.ModerationStatus, c.Ids,
		)
	}
	return nil
}

func (c *Comment) Delete(writer io.Writer) error {
	c.EnsureService()
	for _, id := range c.Ids {
		call := c.Service.Comments.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteComment, err)
		}
		_, _ = fmt.Fprintf(writer, "Comment %s deleted\n", id)
	}
	return nil
}

func WithIds(ids []string) Option {
	return func(c *Comment) {
		c.Ids = ids
	}
}

func WithAuthorChannelId(authorChannelId string) Option {
	return func(c *Comment) {
		c.AuthorChannelId = authorChannelId
	}
}

func WithCanRate(canRate *bool) Option {
	return func(c *Comment) {
		if canRate != nil {
			c.CanRate = canRate
		}
	}
}

func WithChannelId(channelId string) Option {
	return func(c *Comment) {
		c.ChannelId = channelId
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *Comment) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		c.MaxResults = maxResults
	}
}

func WithParentId(parentId string) Option {
	return func(c *Comment) {
		c.ParentId = parentId
	}
}

func WithTextFormat(textFormat string) Option {
	return func(c *Comment) {
		c.TextFormat = textFormat
	}
}

func WithTextOriginal(textOriginal string) Option {
	return func(c *Comment) {
		c.TextOriginal = textOriginal
	}
}

func WithModerationStatus(moderationStatus string) Option {
	return func(c *Comment) {
		c.ModerationStatus = moderationStatus
	}
}

func WithBanAuthor(banAuthor *bool) Option {
	return func(c *Comment) {
		if banAuthor != nil {
			c.BanAuthor = banAuthor
		}
	}
}

func WithVideoId(videoId string) Option {
	return func(c *Comment) {
		c.VideoId = videoId
	}
}

func WithViewerRating(viewerRating string) Option {
	return func(c *Comment) {
		c.ViewerRating = viewerRating
	}
}

var (
	WithParts    = common.WithParts[*Comment]
	WithOutput   = common.WithOutput[*Comment]
	WithJsonpath = common.WithJsonpath[*Comment]
	WithService  = common.WithService[*Comment]
)
