package channelSection

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List channel sections"
	listLong     = "List channel sections"
	listIdsUsage = "Return the channel sections with the given ids"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	channelSectionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", false, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var listTool = mcp.NewTool(
	"channelSection-list",
	mcp.WithTitleAnnotation(listShort),
	mcp.WithDescription(listLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(listIdsUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"hl", mcp.DefaultString(""),
		mcp.Description(hlUsage), mcp.Required(),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(mineUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet"}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(pkg.PartsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "table"),
		mcp.DefaultString("yaml"), mcp.Description(pkg.TableUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(pkg.JPUsage), mcp.Required(),
	),
)

func listHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	channelId, _ = args["channelId"].(string)
	hl, _ = args["hl"].(string)
	mineRaw, _ := args["mine"].(string)
	mine = utils.BoolPtr(mineRaw)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "channelSection list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "channelSection list failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "channelSection list completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
	cs := channelSection.NewChannelSection(
		channelSection.WithIDs(ids),
		channelSection.WithChannelId(channelId),
		channelSection.WithHl(hl),
		channelSection.WithMine(mine),
		channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channelSection.WithService(nil),
	)

	return cs.List(parts, output, jpath, writer)
}
