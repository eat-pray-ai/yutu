// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetLiveStream    = errors.New("failed to get live stream")
	errInsertLiveStream = errors.New("failed to insert live stream")
	errUpdateLiveStream = errors.New("failed to update live stream")
	errDeleteLiveStream = errors.New("failed to delete live stream")
)

type LiveStream struct {
	common.Fields
	Title         string `yaml:"title" json:"title,omitempty"`
	Description   string `yaml:"description" json:"description,omitempty"`
	Mine          *bool  `yaml:"mine" json:"mine,omitempty"`
	FrameRate     string `yaml:"frame_rate" json:"frame_rate,omitempty"`
	IngestionType string `yaml:"ingestion_type" json:"ingestion_type,omitempty"`
	Resolution    string `yaml:"resolution" json:"resolution,omitempty"`

	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type ILiveStream[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
}

type Option func(*LiveStream)

func NewLiveStream(opts ...Option) ILiveStream[youtube.LiveStream] {
	s := &LiveStream{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *LiveStream) Get() ([]*youtube.LiveStream, error) {
	if err := s.EnsureService(); err != nil {
		return nil, err
	}
	call := s.Service.LiveStreams.List(s.Parts)
	if len(s.Ids) > 0 {
		call = call.Id(s.Ids...)
	}
	if s.Mine != nil && *s.Mine {
		call = call.Mine(true)
	}
	if s.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}
	if s.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.OnBehalfOfContentOwnerChannel)
	}

	return common.Paginate(
		&s.Fields, call,
		func(r *youtube.LiveStreamListResponse) ([]*youtube.LiveStream, string) {
			return r.Items, r.NextPageToken
		}, errGetLiveStream,
	)
}

func (s *LiveStream) List(writer io.Writer) error {
	streams, err := s.Get()
	if err != nil && streams == nil {
		return err
	}

	common.PrintList(
		s.Output, streams, writer,
		table.Row{"ID", "Title", "Status"},
		func(stream *youtube.LiveStream) table.Row {
			title := ""
			status := ""
			if stream.Snippet != nil {
				title = stream.Snippet.Title
			}
			if stream.Status != nil {
				status = stream.Status.StreamStatus
			}
			return table.Row{stream.Id, title, status}
		},
	)
	return err
}

func (s *LiveStream) Insert(writer io.Writer) error {
	if err := s.EnsureService(); err != nil {
		return err
	}
	stream := &youtube.LiveStream{
		Snippet: &youtube.LiveStreamSnippet{
			Title:       s.Title,
			Description: s.Description,
		},
		Cdn: &youtube.CdnSettings{
			FrameRate:     s.FrameRate,
			IngestionType: s.IngestionType,
			Resolution:    s.Resolution,
		},
	}

	call := s.Service.LiveStreams.Insert(s.Parts, stream)
	if s.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}
	if s.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertLiveStream, err)
	}

	common.PrintResult(
		s.Output, res, writer, "Live stream inserted: %s\n", res.Id,
	)
	return nil
}

func (s *LiveStream) Update(writer io.Writer) error {
	if err := s.EnsureService(); err != nil {
		return err
	}
	s.Parts = []string{"id", "snippet", "cdn"}
	streams, err := s.Get()
	if err != nil {
		return errors.Join(errUpdateLiveStream, err)
	}
	if len(streams) == 0 {
		return errGetLiveStream
	}

	original := streams[0]
	stream := &youtube.LiveStream{
		Id: original.Id,
		Snippet: &youtube.LiveStreamSnippet{
			Title:       original.Snippet.Title,
			Description: original.Snippet.Description,
		},
		Cdn: &youtube.CdnSettings{},
	}
	if original.Cdn != nil {
		stream.Cdn.FrameRate = original.Cdn.FrameRate
		stream.Cdn.IngestionType = original.Cdn.IngestionType
		stream.Cdn.Resolution = original.Cdn.Resolution
	}

	if s.Title != "" {
		stream.Snippet.Title = s.Title
	}
	if s.Description != "" {
		stream.Snippet.Description = s.Description
	}
	if s.FrameRate != "" {
		stream.Cdn.FrameRate = s.FrameRate
	}
	if s.IngestionType != "" {
		stream.Cdn.IngestionType = s.IngestionType
	}
	if s.Resolution != "" {
		stream.Cdn.Resolution = s.Resolution
	}

	call := s.Service.LiveStreams.Update([]string{"snippet", "cdn"}, stream)
	if s.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}
	if s.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(s.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateLiveStream, err)
	}

	common.PrintResult(
		s.Output, res, writer, "Live stream updated: %s\n", res.Id,
	)
	return nil
}

func (s *LiveStream) Delete(writer io.Writer) error {
	if err := s.EnsureService(); err != nil {
		return err
	}
	for _, id := range s.Ids {
		call := s.Service.LiveStreams.Delete(id)
		if s.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
		}
		if s.OnBehalfOfContentOwnerChannel != "" {
			call = call.OnBehalfOfContentOwnerChannel(s.OnBehalfOfContentOwnerChannel)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteLiveStream, err)
		}

		_, _ = fmt.Fprintf(writer, "Live stream %s deleted\n", id)
	}
	return nil
}

func WithTitle(title string) Option {
	return func(s *LiveStream) {
		s.Title = title
	}
}

func WithDescription(description string) Option {
	return func(s *LiveStream) {
		s.Description = description
	}
}

func WithMine(mine *bool) Option {
	return func(s *LiveStream) {
		if mine != nil {
			s.Mine = mine
		}
	}
}

func WithFrameRate(frameRate string) Option {
	return func(s *LiveStream) {
		s.FrameRate = frameRate
	}
}

func WithIngestionType(ingestionType string) Option {
	return func(s *LiveStream) {
		s.IngestionType = ingestionType
	}
}

func WithResolution(resolution string) Option {
	return func(s *LiveStream) {
		s.Resolution = resolution
	}
}

func WithOnBehalfOfContentOwnerChannel(channel string) Option {
	return func(s *LiveStream) {
		s.OnBehalfOfContentOwnerChannel = channel
	}
}

var (
	WithMaxResults = common.WithMaxResults[*LiveStream]
	WithParts      = common.WithParts[*LiveStream]
	WithOutput     = common.WithOutput[*LiveStream]
	WithService    = common.WithService[*LiveStream]
	WithIds        = common.WithIds[*LiveStream]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*LiveStream]
)
