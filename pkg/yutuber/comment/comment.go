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
	errSetModerationStatus = errors.New("failed to set comment moderation status")
)

type comment struct {
	IDs              []string `yaml:"ids" json:"ids"`
	MaxResults       int64    `yaml:"max_results" json:"max_results"`
	ParentId         string   `yaml:"parent_id" json:"parent_id"`
	TextFormat       string   `yaml:"text_format" json:"text_format"`
	ModerationStatus string   `yaml:"moderation_status" json:"moderation_status"`
	BanAuthor        *bool    `yaml:"ban_author" json:"ban_author"`
}

type Comment interface {
	get([]string) []*youtube.Comment
	List([]string, string)
	Insert(silent bool)
	Update(silent bool)
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

	if c.IDs[0] != "" {
		call = call.Id(c.IDs[0])
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

	return res.Items
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
	// TODO implement me
	panic("implement me")
}

func (c *comment) Update(silent bool) {
	// TODO implement me
	panic("implement me")
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
