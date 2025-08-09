package i18nLanguage

import (
	"bytes"
	"context"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/i18nLanguage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func init() {
	// cmd.MCP.AddTool(listTool, listHandler)
	cmd.MCP.AddResource(langsResource, langsHandler)
	cmd.MCP.AddResource(hlResource, hlHandler)
	i18nLanguageCmd.AddCommand(listCmd)
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

var hlResource = mcp.NewResource(
	hlURI, hlName,
	mcp.WithMIMEType(cmd.JsonMIME),
	mcp.WithResourceDescription(hlDesc),
)

func hlHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	output = "json"
	jpath = "$.*.snippet.hl"
	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		return nil, err
	}

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      hlURI,
			MIMEType: cmd.JsonMIME,
			Text:     writer.String(),
		},
	}
	return contents, nil
}

var langsResource = mcp.NewResource(
	langURI, langName,
	mcp.WithMIMEType(cmd.JsonMIME),
	mcp.WithResourceDescription(long),
)

func langsHandler(
	ctx context.Context, request mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	parts = defaultParts
	output = "json"
	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		return nil, err
	}

	contents := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      langURI,
			MIMEType: cmd.JsonMIME,
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
