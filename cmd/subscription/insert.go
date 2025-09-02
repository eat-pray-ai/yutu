package subscription

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	insertShort    = "Insert a YouTube subscription"
	insertLong     = "Insert a YouTube subscription"
	insertCidUsage = "ID of the channel to be subscribed"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var insertTool = mcp.NewTool(
	"subscription-insert",
	mcp.WithTitleAnnotation(insertShort),
	mcp.WithDescription(insertLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"subscriberChannelId", mcp.DefaultString(""),
		mcp.Description(scidUsage), mcp.Required(),
	),
	mcp.WithString(
		"description", mcp.DefaultString(""),
		mcp.Description(descUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(insertCidUsage), mcp.Required(),
	),
	mcp.WithString(
		"title", mcp.DefaultString(""),
		mcp.Description(titleUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "silent", ""),
		mcp.DefaultString("yaml"), mcp.Description(pkg.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(pkg.JPUsage), mcp.Required(),
	),
)

func insertHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	subscriberChannelId, _ = args["subscriberChannelId"].(string)
	description, _ = args["description"].(string)
	channelId, _ = args["channelId"].(string)
	title, _ = args["title"].(string)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "subscription insert started")

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "subscription insert failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "subscription insert completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func insert(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithSubscriberChannelId(subscriberChannelId),
		subscription.WithDescription(description),
		subscription.WithChannelId(channelId),
		subscription.WithTitle(title),
		subscription.WithService(nil),
	)

	return s.Insert(output, jpath, writer)
}
