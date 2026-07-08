// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package abuseReport

import (
	"encoding/json"
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
	insertShort   = "Insert an abuse report"
	insertLong    = "Insert an abuse report. Use this tool to report abusive content on YouTube such as spam, harassment, or violent content."
	insertExample = `# Report a video as spam
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
			Items: &jsonschema.Schema{Type: "string"}, Default: json.RawMessage(`["snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
		}, cobramcp.GenToolHandler(
			insertTool, func(input abuseReport.AbuseReport, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	abuseReportCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringSliceVarP(
		&abuseTypes, "abuseTypes", "a", []string{}, abuseTypesUsage,
	)
	insertCmd.Flags().StringVarP(
		&description, "description", "d", "", descriptionUsage,
	)
	insertCmd.Flags().StringVarP(&subjectId, "subjectId", "s", "", subjectIdUsage)
	insertCmd.Flags().StringVarP(
		&subjectTypeId, "subjectTypeId", "t", "", subjectTypeUsage,
	)
	insertCmd.Flags().StringVarP(
		&subjectUrl, "subjectUrl", "u", "", subjectUrlUsage,
	)
	insertCmd.Flags().StringVarP(
		&relatedEntityId, "relatedEntityId", "r", "", relatedEntityUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("abuseTypes")
	_ = insertCmd.MarkFlagRequired("subjectId")
	_ = insertCmd.MarkFlagRequired("subjectTypeId")
	cmd.AddMutationFlags(insertCmd)
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would report abuse for %s %s", subjectTypeId, subjectId,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := abuseReport.NewAbuseReport(
			abuseReport.WithAbuseTypes(abuseTypes),
			abuseReport.WithDescription(description),
			abuseReport.WithSubjectId(subjectId),
			abuseReport.WithSubjectTypeId(subjectTypeId),
			abuseReport.WithSubjectUrl(subjectUrl),
			abuseReport.WithRelatedEntityId(relatedEntityId),
			abuseReport.WithParts(parts),
			abuseReport.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}
