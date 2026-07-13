// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetLiveBroadcast            = errors.New("failed to get live broadcast")
	errInsertLiveBroadcast         = errors.New("failed to insert live broadcast")
	errUpdateLiveBroadcast         = errors.New("failed to update live broadcast")
	errDeleteLiveBroadcast         = errors.New("failed to delete live broadcast")
	errBindLiveBroadcast           = errors.New("failed to bind live broadcast")
	errTransitionLiveBroadcast     = errors.New("failed to transition live broadcast")
	errInsertCuepointLiveBroadcast = errors.New("failed to insert cuepoint")
)

type LiveBroadcast struct {
	common.Fields
	Title              string `yaml:"title" json:"title,omitempty"`
	Description        string `yaml:"description" json:"description,omitempty"`
	Mine               *bool  `yaml:"mine" json:"mine,omitempty"`
	BroadcastStatus    string `yaml:"broadcast_status" json:"broadcast_status,omitempty"`
	BroadcastType      string `yaml:"broadcast_type" json:"broadcast_type,omitempty"`
	PrivacyStatus      string `yaml:"privacy_status" json:"privacy_status,omitempty"`
	ScheduledStartTime string `yaml:"scheduled_start_time" json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   string `yaml:"scheduled_end_time" json:"scheduled_end_time,omitempty"`
	StreamId           string `yaml:"stream_id" json:"stream_id,omitempty"`

	CueType              string `yaml:"cue_type" json:"cue_type,omitempty"`
	CueDurationSecs      int64  `yaml:"cue_duration_secs" json:"cue_duration_secs,omitempty"`
	CueInsertionOffsetMs int64  `yaml:"cue_insertion_offset_ms" json:"cue_insertion_offset_ms,omitempty"`
	CueWalltimeMs        uint64 `yaml:"cue_walltime_ms" json:"cue_walltime_ms,omitempty"`

	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type ILiveBroadcast[T any] interface {
	List(io.Writer) error
	Get() ([]*T, error)
	Insert(io.Writer) error
	Update(io.Writer) error
	Delete(io.Writer) error
	Bind(io.Writer) error
	Transition(io.Writer) error
	InsertCuepoint(io.Writer) error
}

type Option func(*LiveBroadcast)

func NewLiveBroadcast(opts ...Option) ILiveBroadcast[youtube.LiveBroadcast] {
	b := &LiveBroadcast{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *LiveBroadcast) Get() ([]*youtube.LiveBroadcast, error) {
	if err := b.EnsureService(); err != nil {
		return nil, err
	}
	call := b.Service.LiveBroadcasts.List(b.Parts)
	if len(b.Ids) > 0 {
		call = call.Id(b.Ids...)
	}
	if b.Mine != nil && *b.Mine {
		call = call.Mine(true)
	}
	if b.BroadcastStatus != "" {
		call = call.BroadcastStatus(b.BroadcastStatus)
	}
	if b.BroadcastType != "" {
		call = call.BroadcastType(b.BroadcastType)
	}
	if b.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
	}
	if b.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
	}

	return common.Paginate(
		&b.Fields, call,
		func(r *youtube.LiveBroadcastListResponse) ([]*youtube.LiveBroadcast, string) {
			return r.Items, r.NextPageToken
		}, errGetLiveBroadcast,
	)
}

func (b *LiveBroadcast) List(writer io.Writer) error {
	broadcasts, err := b.Get()
	if err != nil && broadcasts == nil {
		return err
	}

	common.PrintList(
		b.Output, broadcasts, writer,
		table.Row{"ID", "Title", "Status", "Privacy"},
		func(bc *youtube.LiveBroadcast) table.Row {
			title := ""
			status := ""
			privacy := ""
			if bc.Snippet != nil {
				title = bc.Snippet.Title
			}
			if bc.Status != nil {
				status = bc.Status.LifeCycleStatus
				privacy = bc.Status.PrivacyStatus
			}
			return table.Row{bc.Id, title, status, privacy}
		},
	)
	return err
}

func (b *LiveBroadcast) Insert(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	broadcast := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              b.Title,
			Description:        b.Description,
			ScheduledStartTime: b.ScheduledStartTime,
			ScheduledEndTime:   b.ScheduledEndTime,
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus: b.PrivacyStatus,
		},
	}

	call := b.Service.LiveBroadcasts.Insert(b.Parts, broadcast)
	if b.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
	}
	if b.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errInsertLiveBroadcast, err)
	}

	common.PrintResult(
		b.Output, res, writer, "Live broadcast inserted: %s\n", res.Id,
	)
	return nil
}

func (b *LiveBroadcast) Update(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	b.Parts = []string{"id", "snippet", "status"}
	broadcasts, err := b.Get()
	if err != nil {
		return errors.Join(errUpdateLiveBroadcast, err)
	}
	if len(broadcasts) == 0 {
		return errGetLiveBroadcast
	}

	original := broadcasts[0]
	broadcast := &youtube.LiveBroadcast{
		Id: original.Id,
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              original.Snippet.Title,
			Description:        original.Snippet.Description,
			ScheduledStartTime: original.Snippet.ScheduledStartTime,
			ScheduledEndTime:   original.Snippet.ScheduledEndTime,
		},
		Status: &youtube.LiveBroadcastStatus{},
	}
	if original.Status != nil {
		broadcast.Status.PrivacyStatus = original.Status.PrivacyStatus
	}

	if b.Title != "" {
		broadcast.Snippet.Title = b.Title
	}
	if b.Description != "" {
		broadcast.Snippet.Description = b.Description
	}
	if b.ScheduledStartTime != "" {
		broadcast.Snippet.ScheduledStartTime = b.ScheduledStartTime
	}
	if b.ScheduledEndTime != "" {
		broadcast.Snippet.ScheduledEndTime = b.ScheduledEndTime
	}
	if b.PrivacyStatus != "" {
		broadcast.Status.PrivacyStatus = b.PrivacyStatus
	}

	call := b.Service.LiveBroadcasts.Update(
		[]string{"snippet", "status"}, broadcast,
	)
	if b.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
	}
	if b.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateLiveBroadcast, err)
	}

	common.PrintResult(
		b.Output, res, writer, "Live broadcast updated: %s\n", res.Id,
	)
	return nil
}

func (b *LiveBroadcast) Delete(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	for _, id := range b.Ids {
		call := b.Service.LiveBroadcasts.Delete(id)
		if b.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
		}
		if b.OnBehalfOfContentOwnerChannel != "" {
			call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteLiveBroadcast, err)
		}

		_, _ = fmt.Fprintf(writer, "Live broadcast %s deleted\n", id)
	}
	return nil
}

func (b *LiveBroadcast) Bind(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	for _, id := range b.Ids {
		call := b.Service.LiveBroadcasts.Bind(id, b.Parts)
		if b.StreamId != "" {
			call = call.StreamId(b.StreamId)
		}
		if b.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
		}
		if b.OnBehalfOfContentOwnerChannel != "" {
			call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
		}

		res, err := call.Do()
		if err != nil {
			return errors.Join(errBindLiveBroadcast, err)
		}

		common.PrintResult(
			b.Output, res, writer, "Live broadcast %s bound to stream %s\n", res.Id,
			b.StreamId,
		)
	}
	return nil
}

func (b *LiveBroadcast) Transition(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	for _, id := range b.Ids {
		call := b.Service.LiveBroadcasts.Transition(b.BroadcastStatus, id, b.Parts)
		if b.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
		}
		if b.OnBehalfOfContentOwnerChannel != "" {
			call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
		}

		res, err := call.Do()
		if err != nil {
			return errors.Join(errTransitionLiveBroadcast, err)
		}

		common.PrintResult(
			b.Output, res, writer,
			"Live broadcast %s transitioned to %s\n", id, b.BroadcastStatus,
		)
	}
	return nil
}

func (b *LiveBroadcast) InsertCuepoint(writer io.Writer) error {
	if err := b.EnsureService(); err != nil {
		return err
	}
	cuepoint := &youtube.Cuepoint{
		CueType:               b.CueType,
		DurationSecs:          b.CueDurationSecs,
		InsertionOffsetTimeMs: b.CueInsertionOffsetMs,
		WalltimeMs:            b.CueWalltimeMs,
	}

	for _, id := range b.Ids {
		call := b.Service.LiveBroadcasts.InsertCuepoint(cuepoint).Id(id)
		if b.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(b.OnBehalfOfContentOwner)
		}
		if b.OnBehalfOfContentOwnerChannel != "" {
			call = call.OnBehalfOfContentOwnerChannel(b.OnBehalfOfContentOwnerChannel)
		}

		res, err := call.Do()
		if err != nil {
			return errors.Join(errInsertCuepointLiveBroadcast, err)
		}

		common.PrintResult(
			b.Output, res, writer, "Cuepoint inserted for broadcast %s: %s\n", id,
			res.Id,
		)
	}
	return nil
}

func WithTitle(title string) Option {
	return func(b *LiveBroadcast) {
		b.Title = title
	}
}

func WithDescription(description string) Option {
	return func(b *LiveBroadcast) {
		b.Description = description
	}
}

func WithMine(mine *bool) Option {
	return func(b *LiveBroadcast) {
		if mine != nil {
			b.Mine = mine
		}
	}
}

func WithBroadcastStatus(broadcastStatus string) Option {
	return func(b *LiveBroadcast) {
		b.BroadcastStatus = broadcastStatus
	}
}

func WithBroadcastType(broadcastType string) Option {
	return func(b *LiveBroadcast) {
		b.BroadcastType = broadcastType
	}
}

func WithPrivacyStatus(privacyStatus string) Option {
	return func(b *LiveBroadcast) {
		b.PrivacyStatus = privacyStatus
	}
}

func WithScheduledStartTime(scheduledStartTime string) Option {
	return func(b *LiveBroadcast) {
		b.ScheduledStartTime = scheduledStartTime
	}
}

func WithScheduledEndTime(scheduledEndTime string) Option {
	return func(b *LiveBroadcast) {
		b.ScheduledEndTime = scheduledEndTime
	}
}

func WithStreamId(streamId string) Option {
	return func(b *LiveBroadcast) {
		b.StreamId = streamId
	}
}

func WithCueType(cueType string) Option {
	return func(b *LiveBroadcast) {
		b.CueType = cueType
	}
}

func WithCueDurationSecs(durationSecs int64) Option {
	return func(b *LiveBroadcast) {
		b.CueDurationSecs = durationSecs
	}
}

func WithCueInsertionOffsetMs(offsetMs int64) Option {
	return func(b *LiveBroadcast) {
		b.CueInsertionOffsetMs = offsetMs
	}
}

func WithCueWalltimeMs(walltimeMs uint64) Option {
	return func(b *LiveBroadcast) {
		b.CueWalltimeMs = walltimeMs
	}
}

func WithOnBehalfOfContentOwnerChannel(channel string) Option {
	return func(b *LiveBroadcast) {
		b.OnBehalfOfContentOwnerChannel = channel
	}
}

var (
	WithMaxResults = common.WithMaxResults[*LiveBroadcast]
	WithParts      = common.WithParts[*LiveBroadcast]
	WithOutput     = common.WithOutput[*LiveBroadcast]
	WithService    = common.WithService[*LiveBroadcast]
	WithIds        = common.WithIds[*LiveBroadcast]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*LiveBroadcast]
)
