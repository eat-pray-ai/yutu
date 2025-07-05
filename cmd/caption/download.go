package caption

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	downloadShort   = "Download caption"
	downloadLong    = "Download caption from a video"
	downloadIdUsage = "ID of the caption to download"
)

func init() {
	cmd.MCP.AddTool(downloadTool, downloadHandler)
	captionCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringSliceVarP(
		&ids, "id", "i", []string{}, downloadIdUsage,
	)
	downloadCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", tfmtUsage)
	downloadCmd.Flags().StringVarP(&tlang, "tlang", "l", "", tlangUsage)
	downloadCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	_ = downloadCmd.MarkFlagRequired("id")
	_ = downloadCmd.MarkFlagRequired("file")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: downloadShort,
	Long:  downloadLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := download(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var downloadTool = mcp.NewTool(
	"caption-download",
	mcp.WithTitleAnnotation(downloadShort),
	mcp.WithDescription(downloadLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(downloadIdUsage), mcp.Required(),
	),
	mcp.WithString(
		"file", mcp.DefaultString(""),
		mcp.Description(fileUsage), mcp.Required(),
	),
	mcp.WithString(
		"tfmt", mcp.DefaultString(""),
		mcp.Description(tfmtUsage), mcp.Required(),
	),
	mcp.WithString(
		"tlang", mcp.DefaultString(""),
		mcp.Description(tlangUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOf", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
)

func downloadHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	file, _ = args["file"].(string)
	tfmt, _ = args["tfmt"].(string)
	tlang, _ = args["tlang"].(string)
	onBehalfOf, _ = args["onBehalfOf"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)

	var writer bytes.Buffer
	err := download(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func download(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithIDs(ids),
		caption.WithFile(file),
		caption.WithTfmt(tfmt),
		caption.WithTlang(tlang),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithService(nil),
	)

	return c.Download(writer)
}
