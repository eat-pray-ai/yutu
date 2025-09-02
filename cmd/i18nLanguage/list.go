package i18nLanguage

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/i18nLanguage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.MCP.AddResource(hlResource, hlHandler)
	cmd.MCP.AddResourceTemplate(langsResource, langsHandler)
	i18nLanguageCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", defaultParts, pkg.PartsUsage,
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

var hlResource = mcp.NewResource(
	hlURI, hlName,
	mcp.WithMIMEType(pkg.JsonMIME),
	mcp.WithResourceDescription(hlDesc),
	mcp.WithAnnotations([]mcp.Role{"user", "assistant"}, 0.51),
)

func hlHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	output = "json"
	jpath = "$.*.snippet.hl"

	slog.InfoContext(ctx, "i18nLanguage hl list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nLanguage hl list failed",
			"error", err,
			"uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nLanguage hl list completed successfully",
		"resultSize", writer.Len(),
	)

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      hlURI,
			MIMEType: pkg.JsonMIME,
			Text:     writer.String(),
		},
	}
	return contents, nil
}

var langsResource = mcp.NewResourceTemplate(
	langURI, langName,
	mcp.WithTemplateMIMEType(pkg.JsonMIME),
	mcp.WithTemplateDescription(long),
	mcp.WithTemplateAnnotations([]mcp.Role{"user", "assistant"}, 0.51),
)

func langsHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	hl = utils.ExtractHl(request.Params.URI)
	output = "json"

	slog.InfoContext(ctx, "i18nLanguage list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nLanguage list failed",
			"error", err,
			"uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nLanguage list completed successfully",
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
	i := i18nLanguage.NewI18nLanguage(
		i18nLanguage.WithHl(hl), i18nLanguage.WithService(nil),
	)

	return i.List(parts, output, jpath, writer)
}
