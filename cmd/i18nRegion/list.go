package i18nRegion

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/i18nRegion"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.MCP.AddResourceTemplate(RegionsResource, regionsHandler)
	i18nRegionCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", defaultParts, cmd.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)
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

var RegionsResource = mcp.NewResourceTemplate(
	regionURI, regionName,
	mcp.WithTemplateMIMEType(cmd.JsonMIME),
	mcp.WithTemplateDescription(long),
	mcp.WithTemplateAnnotations([]mcp.Role{"user", "assistant"}, 0.51),
)

func regionsHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	hl = utils.ExtractHl(request.Params.URI)
	output = "json"
	
	slog.InfoContext(ctx, "i18nRegion list started")
	
	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nRegion list failed",
			"error", err,
			"uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nRegion list completed successfully",
		"resultSize", writer.Len(),
	)

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: cmd.JsonMIME,
			Text:     writer.String(),
		},
	}
	return contents, nil
}

func list(writer io.Writer) error {
	i := i18nRegion.NewI18nRegion(
		i18nRegion.WithHl(hl), i18nRegion.WithService(nil),
	)

	return i.List(parts, output, jpath, writer)
}
