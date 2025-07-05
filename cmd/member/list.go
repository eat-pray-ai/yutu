package member

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/member"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	memberCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&memberChannelId, "memberChannelId", "c", "", mcidUsage,
	)
	listCmd.Flags().StringVarP(
		&hasAccessToLevel, "hasAccessToLevel", "a", "", hatlUsage,
	)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().StringVarP(&mode, "mode", "m", "all_current", mmUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var listTool = mcp.NewTool(
	"member-list",
	mcp.WithTitleAnnotation(short),
	mcp.WithDescription(long),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithString(
		"memberChannelId", mcp.DefaultString(""),
		mcp.Description(mcidUsage), mcp.Required(),
	),
	mcp.WithString(
		"hasAccessToLevel", mcp.DefaultString(""),
		mcp.Description(hatlUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(mrUsage), mcp.Required(),
	),
	mcp.WithString(
		"mode", mcp.DefaultString("all_current"),
		mcp.Description(mmUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"snippet"}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(partsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.DefaultString("table"),
		mcp.Description(cmd.TableUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JpUsage), mcp.Required(),
	),
)

func listHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	memberChannelId, _ = args["memberChannelId"].(string)
	hasAccessToLevel, _ = args["hasAccessToLevel"].(string)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	mode, _ = args["mode"].(string)
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
	m := member.NewMember(
		member.WithMemberChannelId(memberChannelId),
		member.WithHasAccessToLevel(hasAccessToLevel),
		member.WithMaxResults(maxResults),
		member.WithMode(mode),
		member.WithService(nil),
	)

	return m.List(parts, output, jpath, writer)
}
