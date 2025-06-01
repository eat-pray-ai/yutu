package commentThread

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service                *youtube.Service
	errGetCommentThread    = errors.New("failed to get comment thread")
	errInsertCommentThread = errors.New("failed to insert comment thread")
)

type commentThread struct {
	ID                           []string `yaml:"id" json:"id"`
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
	get([]string) []*youtube.CommentThread
	List([]string, string)
	Insert(output string)
}

type Option func(*commentThread)

func NewCommentThread(opts ...Option) CommentThread {
	c := &commentThread{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *commentThread) get(parts []string) []*youtube.CommentThread {
	call := service.CommentThreads.List(parts)

	if c.ID != nil {
		call = call.Id(c.ID...)
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
		utils.PrintJSON(c, nil)
		log.Fatalln(errors.Join(errGetCommentThread, err))
	}

	return res.Items
}

func (c *commentThread) List(parts []string, output string) {
	commentThreads := c.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(commentThreads, nil)
	case "yaml":
		utils.PrintYAML(commentThreads, nil)
	default:
		fmt.Println("ID\tTopLevelCommentID")
		for _, commentThread := range commentThreads {
			fmt.Printf(
				"%s\t%s\n", commentThread.Id, commentThread.Snippet.TopLevelComment.Id,
			)
		}
	}
}

func (c *commentThread) Insert(output string) {
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
		utils.PrintJSON(ct, nil)
		log.Fatalln(errors.Join(errInsertCommentThread, err))
	}

	switch output {
	case "json":
		utils.PrintJSON(res, nil)
	case "yaml":
		utils.PrintYAML(res, nil)
	case "silent":
	default:
		fmt.Printf("CommentThread inserted: %s\n", res.Id)
	}
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

func WithID(id []string) Option {
	return func(c *commentThread) {
		c.ID = id
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
