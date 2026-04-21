// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"encoding/json"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "channel-update"
	updateIdUsage = "ID of the channel to update"
	updateShort   = "Update channel information"
	updateLong    = "Update channel information. Use this tool to update channel information."
	updateExample = `# Update channel description
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --description 'New description'
# Update channel title and country
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --title 'New Title' --country US
# Update channel default language
yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --defaultLanguage en`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"country":          {Type: "string", Description: countryUsage},
		"custom_url":       {Type: "string", Description: curlUsage},
		"default_language": {Type: "string", Description: dlUsage},
		"description":      {Type: "string", Description: descUsage},
		"title":            {Type: "string", Description: titleUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			updateTool, func(input channel.Channel, writer io.Writer) error {
				return input.Update(writer)
			},
		),
	)
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&country, "country", "c", "", countryUsage)
	updateCmd.Flags().StringVarP(&customUrl, "customUrl", "u", "", curlUsage)
	updateCmd.Flags().StringVarP(
		&defaultLanguage, "defaultLanguage", "l", "", dlUsage,
	)
	updateCmd.Flags().StringVarP(
		&description, "description", "d", "", descUsage,
	)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)

	_ = updateCmd.MarkFlagRequired("id")
	cmd.AddMutationFlags(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	Run: func(c *cobra.Command, args []string) {
		err := cmd.Confirm(c, "Would update channel: %s", strings.Join(ids, ", "))
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := channel.NewChannel(
			channel.WithIds(ids),
			channel.WithCountry(country),
			channel.WithCustomUrl(customUrl),
			channel.WithDefaultLanguage(defaultLanguage),
			channel.WithDescription(description),
			channel.WithTitle(title),
			channel.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}
