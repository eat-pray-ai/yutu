package videoCategory

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/videoCategory"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.MCP.AddResourceTemplate(categoriesResource, categoriesHandler)
	videoCategoryCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "US", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)
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

var categoriesResource = mcp.NewResourceTemplate(
	vcURI, vcName,
	mcp.WithTemplateMIMEType(pkg.JsonMIME),
	mcp.WithTemplateDescription(long),
	mcp.WithTemplateAnnotations([]mcp.Role{"user", "assistant"}, 0.51),
)

func categoriesHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	hl = utils.ExtractHl(request.Params.URI)
	output = "json"

	slog.InfoContext(ctx, "videoCategory list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "videoCategory list failed",
			"error", err,
			"uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "videoCategory list completed successfully",
		"resultSize", writer.Len(),
	)

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: pkg.JsonMIME,
			Text:     writer.String(),
		},
	}
	return contents, nil
}

func list(writer io.Writer) error {
	vc := videoCategory.NewVideoCategory(
		videoCategory.WithIDs(ids),
		videoCategory.WithHl(hl),
		videoCategory.WithRegionCode(regionCode),
		videoCategory.WithService(nil),
	)

	return vc.List(parts, output, jpath, writer)
}
