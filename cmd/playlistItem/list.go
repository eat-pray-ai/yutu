package playlistItem

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List playlist items"
	listLong     = "List playlist items' info, such as title, description, etc"
	listIdsUsage = "IDs of the playlist items to list"
	listPidUsage = "Return the playlist items within the given playlist"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	playlistItemCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&playlistId, "playlistId", "y", "", listPidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, cmd.MRUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, cmd.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)
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
	"playlistItem-list",
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
		"playlistId", mcp.DefaultString(""),
		mcp.Description(listPidUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoId", mcp.DefaultString(""),
		mcp.Description(vidUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet", "status"}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(cmd.PartsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "table"),
		mcp.DefaultString("yaml"), mcp.Description(cmd.TableUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JPUsage), mcp.Required(),
	),
)

func listHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	playlistId, _ = args["playlistId"].(string)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	videoId, _ = args["videoId"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
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
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithIDs(ids),
		playlistItem.WithPlaylistId(playlistId),
		playlistItem.WithMaxResults(maxResults),
		playlistItem.WithVideoId(videoId),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.List(parts, output, jpath, writer)
}
