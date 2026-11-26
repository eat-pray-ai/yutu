// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package abuseReport

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/abuseReport"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "abuseReport-insert"
	insertConfirm = "Create abuse report for %s %s"
	insertShort   = "Create an abuse report"
	insertLong    = "Create an abuse report. Use this tool to report abusive content on YouTube such as spam, harassment, or violent content."
	insertExample = `  # Report a video as spam
  yutu abuseReport insert --abuseTypes spam --subjectId VIDEO_ID --subjectTypeId video --description "This video is spam"
  # Report a comment as harassment
  yutu abuseReport insert --abuseTypes harassment --subjectId COMMENT_ID --subjectTypeId comment --description "Harassment in comment"`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"abuse_types", "subject_id", "subject_type_id"},
	Properties: map[string]*jsonschema.Schema{
		"abuse_types": {
			Type: "array", Description: abuseTypesUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"description":       {Type: "string", Description: descriptionUsage},
		"subject_id":        {Type: "string", Description: subjectIdUsage},
		"subject_type_id":   {Type: "string", Description: subjectTypeUsage},
		"subject_url":       {Type: "string", Description: subjectUrlUsage},
		"related_entity_id": {Type: "string", Description: relatedEntityUsage},
		"parts": {
			Type: "array", Description: "Parts to include in the response",
			Items: &jsonschema.Schema{Type: "string"}, Default: jsontext.Value(`["snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input abuseReport.AbuseReport) string {
					return fmt.Sprintf(insertConfirm, input.SubjectTypeId, input.SubjectId)
				},
				func(input abuseReport.AbuseReport, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	abuseReportCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	abuseTypes      []string
	description     string
	subjectId       string
	subjectTypeId   string
	subjectUrl      string
	relatedEntityId string
	parts           []string
	output          string
	yes             bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(
		&o.abuseTypes, "abuseTypes", "a", nil, abuseTypesUsage,
	)
	c.Flags().StringVarP(
		&o.description, "description", "d", "", descriptionUsage,
	)
	c.Flags().StringVarP(&o.subjectId, "subjectId", "s", "", subjectIdUsage)
	c.Flags().StringVarP(
		&o.subjectTypeId, "subjectTypeId", "t", "", subjectTypeUsage,
	)
	c.Flags().StringVarP(
		&o.subjectUrl, "subjectUrl", "u", "", subjectUrlUsage,
	)
	c.Flags().StringVarP(
		&o.relatedEntityId, "relatedEntityId", "r", "", relatedEntityUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("abuseTypes")
	_ = c.MarkFlagRequired("subjectId")
	_ = c.MarkFlagRequired("subjectTypeId")
}

func (o *insertOpts) build() *abuseReport.AbuseReport {
	return abuseReport.NewAbuseReport(
		abuseReport.WithAbuseTypes(o.abuseTypes),
		abuseReport.WithDescription(o.description),
		abuseReport.WithSubjectId(o.subjectId),
		abuseReport.WithSubjectTypeId(o.subjectTypeId),
		abuseReport.WithSubjectUrl(o.subjectUrl),
		abuseReport.WithRelatedEntityId(o.relatedEntityId),
		abuseReport.WithParts(o.parts),
		abuseReport.WithOutput(o.output),
	).(*abuseReport.AbuseReport)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.subjectTypeId, opts.subjectId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}
