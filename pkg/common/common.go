// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"fmt"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/youtube/v3"
)

type Fields struct {
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

func (d *Fields) EnsureService() error {
	if d.Service == nil {
		svc, err := auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
		if err != nil {
			return fmt.Errorf("failed to create YouTube service: %w", err)
		}
		d.Service = svc
	}
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
	pageToken := ""
	for f.MaxResults > 0 {
		call = call.MaxResults(min(f.MaxResults, pkg.PerPage))
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errWrap, err)
		}
		got, nextToken := extract(res)
		f.MaxResults -= pkg.PerPage
		items = append(items, got...)
		pageToken = nextToken
		if pageToken == "" || len(got) == 0 {
			break
		}
	}
	return items, nil
}
