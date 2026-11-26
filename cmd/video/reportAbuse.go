// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	reportAbuseTool    = "video-reportAbuse"
	raIdsUsage         = "IDs of the videos to report abuse on"
	raLangUsage        = "Language that the content was viewed in"
	reportAbuseConfirm = "Report abuse on video(s) %s"
	reportAbuseShort   = "Report abuse on a video"
	reportAbuseLong    = "Report abuse on a video. Use this tool to report abuse on a video."
	reportAbuseExample = `  # Report abuse on a video
  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V
  # Report abuse with secondary reason and language
  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V --secondaryReasonId V1 --language en
  # Report abuse with comments
  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId N --comments 'Spam content'`
)

var reportAbuseInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "reason_id"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: raIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"reason_id":                  {Type: "string", Description: ridUsage},
		"secondary_reason_id":        {Type: "string", Description: sridUsage},
		"comments":                   {Type: "string", Description: commentsUsage},
		"language":                   {Type: "string", Description: raLangUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: reportAbuseTool, Title: reportAbuseShort,
			Description: reportAbuseLong,
			InputSchema: reportAbuseInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			reportAbuseTool, cobramcp.ConfirmThen(
				func(input video.Video) string {
					return fmt.Sprintf(reportAbuseConfirm, strings.Join(input.Ids, ", "))
				},
				func(input video.Video, w io.Writer) error {
					return input.ReportAbuse(w)
				},
			),
		),
	)
	videoCmd.AddCommand(newReportAbuseCmd())
}

type reportAbuseOpts struct {
	ids               []string
	reasonId          string
	secondaryReasonId string
	comments          string
	language          string
	oboco             string
	yes               bool
}

func (o *reportAbuseOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, raIdsUsage)
	c.Flags().StringVarP(&o.reasonId, "reasonId", "r", "", ridUsage)
	c.Flags().StringVarP(
		&o.secondaryReasonId, "secondaryReasonId", "s", "", sridUsage,
	)
	c.Flags().StringVarP(&o.comments, "comments", "c", "", commentsUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", raLangUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
	_ = c.MarkFlagRequired("reasonId")
}

func (o *reportAbuseOpts) build(f *pflag.FlagSet) *video.Video {
	return video.NewVideo(
		video.WithIds(o.ids),
		video.WithReasonId(o.reasonId),
		video.WithSecondaryReasonId(o.secondaryReasonId),
		video.WithComments(o.comments),
		utils.OptIfChanged(f, "language", &o.language, video.WithLanguage),
		video.WithOBOCO(o.oboco),
	).(*video.Video)
}

func newReportAbuseCmd() *cobra.Command {
	var opts reportAbuseOpts

	c := &cobra.Command{
		Use:     "reportAbuse",
		Short:   reportAbuseShort,
		Long:    reportAbuseLong,
		Example: reportAbuseExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(reportAbuseConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).ReportAbuse(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}
