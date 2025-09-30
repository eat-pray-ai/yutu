// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

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

type Comment[T any] interface {
	Get([]string) ([]*T, error)
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Delete(io.Writer) error
	MarkAsSpam(string, string, io.Writer) error
	SetModerationStatus(string, string, io.Writer) error
}

type Option func(*comment)

func NewComment(opts ...Option) Comment[youtube.Comment] {
	c := &comment{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *comment) Get(parts []string) ([]*youtube.Comment, error) {
	call := service.Comments.List(parts)
	if c.IDs[0] != "" {
		call = call.Id(c.IDs...)
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

func (c *comment) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	comments, err := c.Get(parts)
	if err != nil && comments == nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(comments, jpath, writer)
	case "yaml":
		utils.PrintYAML(comments, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
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

func (c *comment) Insert(output string, jpath string, writer io.Writer) error {
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

	call := service.Comments.Insert([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertComment, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment inserted: %s\n", res.Id)
	}
	return nil
}

func (c *comment) Update(output string, jpath string, writer io.Writer) error {
	comments, err := c.Get([]string{"id", "snippet"})
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

	call := service.Comments.Update([]string{"snippet"}, comment)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateComment, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment updated: %s\n", res.Id)
	}
	return nil
}

func (c *comment) MarkAsSpam(
	output string, jpath string, writer io.Writer,
) error {
	call := service.Comments.MarkAsSpam(c.IDs)
	err := call.Do()
	if err != nil {
		return errors.Join(errMarkAsSpam, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(c, jpath, writer)
	case "yaml":
		utils.PrintYAML(c, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Comment marked as spam: %s\n", c.IDs)
	}
	return nil
}

func (c *comment) SetModerationStatus(
	output string, jpath string, writer io.Writer,
) error {
	call := service.Comments.SetModerationStatus(c.IDs, c.ModerationStatus)

	if c.BanAuthor != nil {
		call = call.BanAuthor(*c.BanAuthor)
	}

	err := call.Do()
	if err != nil {
		return errors.Join(errSetModerationStatus, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(c, jpath, writer)
	case "yaml":
		utils.PrintYAML(c, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(
			writer, "Comment moderation status set to %s: %s\n",
			c.ModerationStatus, c.IDs,
		)
	}
	return nil
}

func (c *comment) Delete(writer io.Writer) error {
	for _, id := range c.IDs {
		call := service.Comments.Delete(id)
		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteComment, err)
		}
		_, _ = fmt.Fprintf(writer, "Comment %s deleted\n", id)
	}
	return nil
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

func WithCanRate(canRate *bool) Option {
	return func(c *comment) {
		if canRate != nil {
			c.CanRate = canRate
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
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
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

func WithBanAuthor(banAuthor *bool) Option {
	return func(c *comment) {
		if banAuthor != nil {
			c.BanAuthor = banAuthor
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

func WithService(svc *youtube.Service) Option {
	return func(_ *comment) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
