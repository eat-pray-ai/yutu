package channel

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List channel's info"
	listLong     = "List channel's info, such as title, description, etc."
	listIdsUsage = "Return the channels with the specified IDs"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&categoryId, "categoryId", "g", "", cidUsage,
	)
	listCmd.Flags().StringVarP(
		&forHandle, "forHandle", "d", "", fhUsage,
	)
	listCmd.Flags().StringVarP(
		&forUsername, "forUsername", "u", "", fuUsage,
	)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().BoolVarP(
		managedByMe, "managedByMe", "E", false, mbmUsage,
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, cmd.MRUsage,
	)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		mySubscribers, "mySubscribers", "S", false, msUsage,
	)
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
	"channel-list",
	mcp.WithTitleAnnotation(listShort),
	mcp.WithDescription(listLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithString(
		"categoryId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"forHandle", mcp.DefaultString(""),
		mcp.Description(fhUsage), mcp.Required(),
	),
	mcp.WithString(
		"forUsername", mcp.DefaultString(""),
		mcp.Description(fuUsage), mcp.Required(),
	),
	mcp.WithString(
		"hl", mcp.DefaultString(""),
		mcp.Description(hlUsage), mcp.Required(),
	),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(listIdsUsage), mcp.Required(),
	),
	mcp.WithString(
		"managedByMe", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(mbmUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(mineUsage), mcp.Required(),
	),
	mcp.WithString(
		"mySubscribers", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(msUsage), mcp.Required(),
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
	categoryId, _ = args["categoryId"].(string)
	forHandle, _ = args["forHandle"].(string)
	forUsername, _ = args["forUsername"].(string)
	hl, _ = args["hl"].(string)
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	managedByMeRaw, _ := args["managedByMe"].(string)
	managedByMe = utils.BoolPtr(managedByMeRaw)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	mineRaw, ok := args["mine"].(string)
	if !ok {
		mineRaw = "true"
	}
	mine = utils.BoolPtr(mineRaw)
	mySubscribersRaw, _ := args["mySubscribers"].(string)
	mySubscribers = utils.BoolPtr(mySubscribersRaw)
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
	c := channel.NewChannel(
		channel.WithCategoryId(categoryId),
		channel.WithForHandle(forHandle),
		channel.WithForUsername(forUsername),
		channel.WithHl(hl),
		channel.WithIDs(ids),
		channel.WithChannelManagedByMe(managedByMe),
		channel.WithMaxResults(maxResults),
		channel.WithMine(mine),
		channel.WithMySubscribers(mySubscribers),
		channel.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channel.WithService(nil),
	)

	return c.List(parts, output, jpath, writer)
}
