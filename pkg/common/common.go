// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// RedirectURL is the OAuth redirect URL used by file-based auth (CLI / stdio mode).
// Set by the command layer before EnsureService is called.
var RedirectURL = "http://localhost:8216"

type Fields struct {
	Ctx        context.Context  `yaml:"-" json:"-"`
	Service    *youtube.Service `yaml:"-" json:"-"`
	Ids        []string         `yaml:"ids" json:"ids,omitempty"`
	MaxResults int64            `yaml:"max_results" json:"max_results,omitempty"`
	Hl         string           `yaml:"hl" json:"hl,omitempty"`
	ChannelId  string           `yaml:"channel_id" json:"channel_id,omitempty"`
	Parts      []string         `yaml:"parts" json:"parts,omitempty"`
	Output     string           `yaml:"output" json:"output,omitempty"`

	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner,omitempty"`
}

func (d *Fields) GetFields() *Fields {
	return d
}

// SetContext implements the cobra-mcp ContextAware interface, enabling
// automatic context injection from MCP tool handlers.
func (d *Fields) SetContext(ctx context.Context) {
	if d != nil {
		d.Ctx = ctx
	}
}

func (d *Fields) EnsureService() error {
	if d.Service != nil {
		return nil
	}

	// MCP OAuth path: use access token from auth middleware context.
	if d.Ctx != nil {
		if tokenInfo := sdkauth.TokenInfoFromContext(d.Ctx); tokenInfo != nil {
			if rawToken, ok := tokenInfo.Extra["access_token"].(string); ok {
				ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: rawToken})
				client := oauth2.NewClient(d.Ctx, ts)
				svc, err := youtube.NewService(d.Ctx, option.WithHTTPClient(client))
				if err != nil {
					return fmt.Errorf("failed to create YouTube service: %w", err)
				}
				d.Service = svc
				return nil
			}
		}
	}

	// File-based auth path (CLI / stdio mode).
	svc, err := auth.NewY2BService(
		auth.WithCredential("", pkg.Root.FS()),
		auth.WithCacheToken("", pkg.Root.FS()),
		auth.WithRedirectURL(RedirectURL),
	).GetService()
	if err != nil {
		return fmt.Errorf("failed to create YouTube service: %w", err)
	}
	d.Service = svc
	return nil
}

type HasFields interface {
	GetFields() *Fields
	EnsureService() error
}

func WithParts[T HasFields](parts []string) func(T) {
	return func(t T) {
		t.GetFields().Parts = parts
	}
}

func WithOutput[T HasFields](output string) func(T) {
	return func(t T) {
		t.GetFields().Output = output
	}
}

func WithService[T HasFields](svc *youtube.Service) func(T) {
	return func(t T) {
		t.GetFields().Service = svc
	}
}

func WithIds[T HasFields](ids []string) func(T) {
	return func(t T) {
		t.GetFields().Ids = ids
	}
}

func WithMaxResults[T HasFields](maxResults int64) func(T) {
	return func(t T) {
		if maxResults < 0 {
			t.GetFields().MaxResults = 1
		} else if maxResults == 0 {
			t.GetFields().MaxResults = math.MaxInt64
		} else {
			t.GetFields().MaxResults = maxResults
		}
	}
}

func WithHl[T HasFields](hl string) func(T) {
	return func(t T) {
		t.GetFields().Hl = hl
	}
}

func WithChannelId[T HasFields](channelId string) func(T) {
	return func(t T) {
		t.GetFields().ChannelId = channelId
	}
}

func WithOnBehalfOfContentOwner[T HasFields](owner string) func(T) {
	return func(t T) {
		t.GetFields().OnBehalfOfContentOwner = owner
	}
}

// PrintList handles the json/yaml/table output switch for List methods.
// The header and row function are only used for table output.
func PrintList[T any](output string, items []*T, w io.Writer, header table.Row, row func(*T) table.Row) {
	switch output {
	case "json":
		utils.PrintJSON(items, w)
	case "yaml":
		utils.PrintYAML(items, w)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(w)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(header)
		for _, item := range items {
			tb.AppendRow(row(item))
		}
	}
}

// PrintResult handles the json/yaml/silent/default output switch for mutation methods.
func PrintResult(output string, data any, w io.Writer, format string, args ...any) {
	switch output {
	case "json":
		utils.PrintJSON(data, w)
	case "yaml":
		utils.PrintYAML(data, w)
	case "silent":
	default:
		_, _ = fmt.Fprintf(w, format, args...)
	}
}

// PagedLister is satisfied by all YouTube API *XxxListCall types.
type PagedLister[C any, R any] interface {
	MaxResults(int64) C
	PageToken(string) C
	Do(opts ...googleapi.CallOption) (*R, error)
}

// Paginate fetches all pages of results. It handles MaxResults, PageToken,
// Do(), and error wrapping automatically. The extract function pulls items
// and the next page token from the response.
func Paginate[C PagedLister[C, R], R any, T any](
	f *Fields, call C,
	extract func(*R) ([]*T, string),
	errWrap error,
) ([]*T, error) {
	var items []*T
	remaining := f.MaxResults
	pageToken := ""
	for remaining > 0 {
		call = call.MaxResults(min(remaining, pkg.PerPage))
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errWrap, err)
		}
		got, nextToken := extract(res)
		remaining -= pkg.PerPage
		items = append(items, got...)
		pageToken = nextToken
		if pageToken == "" || len(got) == 0 {
			break
		}
	}
	return items, nil
}
