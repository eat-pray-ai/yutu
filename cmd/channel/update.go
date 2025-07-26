package channel

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	updateShort   = "Update channel's info"
	updateLong    = "Update channel's info, such as title, description, etc"
	updateIdUsage = "ID of the channel to update"
)

func init() {
	cmd.MCP.AddTool(updateTool, updateHandler)
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
	updateCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var updateTool = mcp.NewTool(
	"channel-update",
	mcp.WithTitleAnnotation(updateShort),
	mcp.WithDescription(updateLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(updateIdUsage), mcp.Required(),
	),
	mcp.WithString(
		"country", mcp.DefaultString(""),
		mcp.Description(countryUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"customUrl", mcp.DefaultString(""),
		mcp.Description(curlUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"defaultLanguage", mcp.DefaultString(""),
		mcp.Description(dlUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"description", mcp.DefaultString(""),
		mcp.Description(descUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"title", mcp.DefaultString(""),
		mcp.Description(titleUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "silent", ""),
		mcp.DefaultString("yaml"), mcp.Description(cmd.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JPUsage),
		mcp.Required(),
	),
)

func updateHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	country, _ = args["country"].(string)
	customUrl, _ = args["customUrl"].(string)
	defaultLanguage, _ = args["defaultLanguage"].(string)
	description, _ = args["description"].(string)
	title, _ = args["title"].(string)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func update(writer io.Writer) error {
	c := channel.NewChannel(
		channel.WithIDs(ids),
		channel.WithCountry(country),
		channel.WithCustomUrl(customUrl),
		channel.WithDefaultLanguage(defaultLanguage),
		channel.WithDescription(description),
		channel.WithTitle(title),
		channel.WithService(nil),
	)

	return c.Update(output, jpath, writer)
}
