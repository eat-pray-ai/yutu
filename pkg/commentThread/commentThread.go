package commentThread

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service                *youtube.Service
	errGetCommentThread    = errors.New("failed to get comment thread")
	errInsertCommentThread = errors.New("failed to insert comment thread")
)

type commentThread struct {
	IDs                          []string `yaml:"ids" json:"ids"`
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

type CommentThread interface {
	Get([]string) ([]*youtube.CommentThread, error)
	List([]string, string, io.Writer) error
	Insert(output string, writer io.Writer) error
}

type Option func(*commentThread)

func NewCommentThread(opts ...Option) CommentThread {
	c := &commentThread{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *commentThread) Get(parts []string) ([]*youtube.CommentThread, error) {
	call := service.CommentThreads.List(parts)

	if len(c.IDs) > 0 {
		call = call.Id(c.IDs...)
	}

	if c.AllThreadsRelatedToChannelId != "" {
		call = call.AllThreadsRelatedToChannelId(c.AllThreadsRelatedToChannelId)
	}

	if c.ChannelId != "" {
		call = call.ChannelId(c.ChannelId)
	}

	if c.MaxResults <= 0 {
		c.MaxResults = 1
	}
	call = call.MaxResults(c.MaxResults)

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

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetCommentThread, err)
	}

	return res.Items, nil
}

func (c *commentThread) List(
	parts []string, output string, writer io.Writer,
) error {
	commentThreads, err := c.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(commentThreads, writer)
	case "yaml":
		utils.PrintYAML(commentThreads, writer)
	default:
		_, _ = fmt.Fprintln(writer, "ID\tTopLevelCommentID")
		for _, commentThread := range commentThreads {
			_, _ = fmt.Fprintf(
				writer, "%s\t%s\n",
				commentThread.Id, commentThread.Snippet.TopLevelComment.Id,
			)
		}
	}
	return nil
}

func (c *commentThread) Insert(output string, writer io.Writer) error {
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
		utils.PrintJSON(res, writer)
	case "yaml":
		utils.PrintYAML(res, writer)
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

func WithIDs(ids []string) Option {
	return func(c *commentThread) {
		c.IDs = ids
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(c *commentThread) {
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
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
