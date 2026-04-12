// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideo    = errors.New("failed to get video")
	errInsertVideo = errors.New("failed to insert video")
	errUpdateVideo = errors.New("failed to update video")
	errRating      = errors.New("failed to rate video")
	errGetRating   = errors.New("failed to get rating")
	errDeleteVideo = errors.New("failed to delete video")
	errReportAbuse = errors.New("failed to report abuse")
)

type Video struct {
	*common.Fields
	AutoLevels  *bool    `yaml:"auto_levels" json:"auto_levels,omitempty"`
	File        string   `yaml:"file" json:"file,omitempty"`
	Title       string   `yaml:"title" json:"title,omitempty"`
	Description string   `yaml:"description" json:"description,omitempty"`
	Tags        []string `yaml:"tags" json:"tags,omitempty"`
	Language    string   `yaml:"language" json:"language,omitempty"`
	Locale      string   `yaml:"locale" json:"locale,omitempty"`
	License     string   `yaml:"license" json:"license,omitempty"`
	Thumbnail   string   `yaml:"thumbnail" json:"thumbnail,omitempty"`
	Rating      string   `yaml:"rating" json:"rating,omitempty"`
	Chart       string   `yaml:"chart" json:"chart,omitempty"`
	Comments    string   `yaml:"comments" json:"comments,omitempty"`
	PlaylistId  string   `yaml:"playlist_id" json:"playlist_id,omitempty"`
	CategoryId  string   `yaml:"category_id" json:"category_id,omitempty"`
	Privacy     string   `yaml:"privacy" json:"privacy,omitempty"`
	ForKids     *bool    `yaml:"for_kids" json:"for_kids,omitempty"`
	Embeddable  *bool    `yaml:"embeddable" json:"embeddable,omitempty"`
	PublishAt   string   `yaml:"publish_at" json:"publish_at,omitempty"`
	RegionCode  string   `yaml:"region_code" json:"region_code,omitempty"`
	ReasonId    string   `yaml:"reason_id" json:"reason_id,omitempty"`
	Stabilize   *bool    `yaml:"stabilize" json:"stabilize,omitempty"`
	MaxHeight   int64    `yaml:"max_height" json:"max_height,omitempty"`
	MaxWidth    int64    `yaml:"max_width" json:"max_width,omitempty"`

	SecondaryReasonId             string `yaml:"secondary_reason_id" json:"secondary_reason_id,omitempty"`
	NotifySubscribers             *bool  `yaml:"notify_subscribers" json:"notify_subscribers,omitempty"`
	PublicStatsViewable           *bool  `yaml:"public_stats_viewable" json:"public_stats_viewable,omitempty"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel,omitempty"`
}

type IVideo[T any] interface {
	List(io.Writer) error
	Insert(io.Writer) error
	Update(io.Writer) error
	Rate(io.Writer) error
	GetRating(io.Writer) error
	Delete(io.Writer) error
	ReportAbuse(io.Writer) error
	Get() ([]*T, error)
}

type Option func(*Video)

func NewVideo(opts ...Option) IVideo[youtube.Video] {
	v := &Video{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *Video) Get() ([]*youtube.Video, error) {
	if err := v.EnsureService(); err != nil {
		return nil, err
	}
	call := v.Service.Videos.List(v.Parts)
	if len(v.Ids) > 0 {
		call = call.Id(v.Ids...)
	}
	if v.Chart != "" {
		call = call.Chart(v.Chart)
	}
	if v.Hl != "" {
		call = call.Hl(v.Hl)
	}
	if v.Locale != "" {
		call = call.Locale(v.Locale)
	}
	if v.CategoryId != "" {
		call = call.VideoCategoryId(v.CategoryId)
	}
	if v.Rating != "" {
		call = call.MyRating(v.Rating)
	}
	if v.RegionCode != "" {
		call = call.RegionCode(v.RegionCode)
	}
	if v.MaxHeight != 0 {
		call = call.MaxHeight(v.MaxHeight)
	}
	if v.MaxWidth != 0 {
		call = call.MaxWidth(v.MaxWidth)
	}
	if v.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
	}

	return common.Paginate(
		v.Fields, call,
		func(r *youtube.VideoListResponse) ([]*youtube.Video, string) {
			return r.Items, r.NextPageToken
		}, errGetVideo,
	)
}

func (v *Video) List(writer io.Writer) error {
	videos, err := v.Get()
	if err != nil && videos == nil {
		return err
	}

	common.PrintList(
		v.Output, videos, writer, table.Row{"ID", "Title", "Channel ID", "Views"},
		func(video *youtube.Video) table.Row {
			title := ""
			channelId := ""
			var views uint64
			if video.Snippet != nil {
				title = video.Snippet.Title
				channelId = video.Snippet.ChannelId
			}
			if video.Statistics != nil {
				views = video.Statistics.ViewCount
			}
			return table.Row{video.Id, title, channelId, views}
		},
	)
	return err
}

func (v *Video) Insert(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	file, err := pkg.Root.Open(v.File)
	if err != nil {
		return errors.Join(errInsertVideo, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if !slices.Contains(v.Tags, "yutu🐰") {
		v.Tags = append(v.Tags, "yutu🐰")
	}

	if v.Title == "" {
		v.Title = utils.GetFileName(v.File)
	}

	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:                v.Title,
			Description:          v.Description,
			Tags:                 v.Tags,
			CategoryId:           v.CategoryId,
			ChannelId:            v.ChannelId,
			DefaultLanguage:      v.Language,
			DefaultAudioLanguage: v.Language,
		},
		Status: &youtube.VideoStatus{
			License:         v.License,
			PublishAt:       v.PublishAt,
			PrivacyStatus:   v.Privacy,
			ForceSendFields: []string{"SelfDeclaredMadeForKids"},
		},
	}

	if v.Embeddable != nil {
		video.Status.Embeddable = *v.Embeddable
	}
	if v.ForKids != nil {
		video.Status.SelfDeclaredMadeForKids = *v.ForKids
	}
	if v.PublicStatsViewable != nil {
		video.Status.PublicStatsViewable = *v.PublicStatsViewable
	}

	call := v.Service.Videos.Insert([]string{"snippet,status"}, video)

	if v.AutoLevels != nil {
		call = call.AutoLevels(*v.AutoLevels)
	}
	if v.NotifySubscribers != nil {
		call = call.NotifySubscribers(*v.NotifySubscribers)
	}
	if v.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
	}
	if v.OnBehalfOfContentOwnerChannel != "" {
		call = call.OnBehalfOfContentOwnerChannel(v.OnBehalfOfContentOwnerChannel)
	}
	if v.Stabilize != nil {
		call = call.Stabilize(*v.Stabilize)
	}

	res, err := call.Media(file).Do()
	if err != nil {
		return errors.Join(errInsertVideo, err)
	}

	if v.Thumbnail != "" {
		t := thumbnail.NewThumbnail(
			thumbnail.WithVideoId(res.Id),
			thumbnail.WithFile(v.Thumbnail),
			thumbnail.WithService(v.Service),
			thumbnail.WithOutput("silent"),
		)
		_ = t.Set(nil)
	}

	if v.PlaylistId != "" {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(res.Snippet.Title),
			playlistItem.WithDescription(res.Snippet.Description),
			playlistItem.WithKind("video"),
			playlistItem.WithKVideoId(res.Id),
			playlistItem.WithPlaylistId(v.PlaylistId),
			playlistItem.WithChannelId(res.Snippet.ChannelId),
			playlistItem.WithPrivacy(res.Status.PrivacyStatus),
			playlistItem.WithService(v.Service),
			playlistItem.WithOutput("silent"),
		)

		_ = pi.Insert(writer)
	}

	common.PrintResult(v.Output, res, writer, "Video inserted: %s\n", res.Id)
	return nil
}

func (v *Video) Update(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	v.Parts = []string{"id", "snippet", "status"}
	videos, err := v.Get()

	if err != nil {
		return errors.Join(errUpdateVideo, err)
	}
	if len(videos) == 0 {
		return errGetVideo
	}

	original := videos[0]

	// Build a new video with only writable fields to avoid sending
	// read-only fields (thumbnails, channelId, etc.) that cause
	// invalidVideoMetadata errors.
	video := &youtube.Video{
		Id:      original.Id,
		Snippet: &youtube.VideoSnippet{},
		Status:  &youtube.VideoStatus{},
	}
	if original.Snippet != nil {
		video.Snippet.Title = original.Snippet.Title
		video.Snippet.Description = original.Snippet.Description
		video.Snippet.Tags = original.Snippet.Tags
		video.Snippet.CategoryId = original.Snippet.CategoryId
		video.Snippet.DefaultLanguage = original.Snippet.DefaultLanguage
	}
	if original.Status != nil {
		video.Status.Embeddable = original.Status.Embeddable
		video.Status.License = original.Status.License
		video.Status.PrivacyStatus = original.Status.PrivacyStatus
		video.Status.PublicStatsViewable = original.Status.PublicStatsViewable
		video.Status.PublishAt = original.Status.PublishAt
		video.Status.SelfDeclaredMadeForKids = original.Status.SelfDeclaredMadeForKids
	}

	if v.Title != "" {
		video.Snippet.Title = v.Title
	}
	if v.Description != "" {
		video.Snippet.Description = v.Description
	}
	if v.Tags != nil {
		if !slices.Contains(v.Tags, "yutu🐰") {
			v.Tags = append(v.Tags, "yutu🐰")
		}
		video.Snippet.Tags = v.Tags
	}
	if v.Language != "" {
		video.Snippet.DefaultLanguage = v.Language
	}
	if v.License != "" {
		video.Status.License = v.License
	}
	if v.CategoryId != "" {
		video.Snippet.CategoryId = v.CategoryId
	}
	if v.Privacy != "" {
		video.Status.PrivacyStatus = v.Privacy
	}
	if v.Embeddable != nil {
		video.Status.Embeddable = *v.Embeddable
	}

	call := v.Service.Videos.Update([]string{"snippet,status"}, video)
	if v.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
	}

	res, err := call.Do()
	if err != nil {
		return errors.Join(errUpdateVideo, err)
	}

	if v.Thumbnail != "" {
		t := thumbnail.NewThumbnail(
			thumbnail.WithVideoId(res.Id),
			thumbnail.WithFile(v.Thumbnail),
			thumbnail.WithService(v.Service),
			thumbnail.WithOutput("silent"),
		)
		_ = t.Set(nil)
	}

	if v.PlaylistId != "" {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(res.Snippet.Title),
			playlistItem.WithDescription(res.Snippet.Description),
			playlistItem.WithKind("video"),
			playlistItem.WithKVideoId(res.Id),
			playlistItem.WithPlaylistId(v.PlaylistId),
			playlistItem.WithChannelId(res.Snippet.ChannelId),
			playlistItem.WithPrivacy(res.Status.PrivacyStatus),
			playlistItem.WithService(v.Service),
			playlistItem.WithOutput("silent"),
		)

		_ = pi.Insert(writer)
	}

	common.PrintResult(v.Output, res, writer, "Video updated: %s\n", res.Id)
	return nil
}

func (v *Video) Rate(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	for _, id := range v.Ids {
		call := v.Service.Videos.Rate(id, v.Rating)
		err := call.Do()
		if err != nil {
			return errors.Join(errRating, err)
		}
		_, _ = fmt.Fprintf(writer, "Video %s rated %s\n", id, v.Rating)
	}
	return nil
}

func (v *Video) GetRating(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	call := v.Service.Videos.GetRating(v.Ids)
	if v.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
	}
	res, err := call.Do()
	if err != nil {
		return errors.Join(errGetRating, err)
	}

	switch v.Output {
	case "json":
		utils.PrintJSON(res.Items, writer)
	case "yaml":
		utils.PrintYAML(res.Items, writer)
	default:
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"ID", "Rating"})
		for _, item := range res.Items {
			tb.AppendRow(table.Row{item.VideoId, item.Rating})
		}
	}
	return nil
}

func (v *Video) Delete(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	for _, id := range v.Ids {
		call := v.Service.Videos.Delete(id)
		if v.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errDeleteVideo, err)
		}
		_, _ = fmt.Fprintf(writer, "Video %s deleted", id)
	}
	return nil
}

func (v *Video) ReportAbuse(writer io.Writer) error {
	if err := v.EnsureService(); err != nil {
		return err
	}
	for _, id := range v.Ids {
		videoAbuseReport := &youtube.VideoAbuseReport{
			Comments:          v.Comments,
			Language:          v.Language,
			ReasonId:          v.ReasonId,
			SecondaryReasonId: v.SecondaryReasonId,
			VideoId:           id,
		}

		call := v.Service.Videos.ReportAbuse(videoAbuseReport)
		if v.OnBehalfOfContentOwner != "" {
			call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwner)
		}

		err := call.Do()
		if err != nil {
			return errors.Join(errReportAbuse, err)
		}

		_, _ = fmt.Fprintf(writer, "Video %s reported for abuse", id)
	}
	return nil
}

func WithAutoLevels(autoLevels *bool) Option {
	return func(v *Video) {
		if autoLevels != nil {
			v.AutoLevels = autoLevels
		}
	}
}

func WithFile(file string) Option {
	return func(v *Video) {
		v.File = file
	}
}

func WithTitle(title string) Option {
	return func(v *Video) {
		v.Title = title
	}
}

func WithDescription(description string) Option {
	return func(v *Video) {
		v.Description = description
	}
}

func WithTags(tags []string) Option {
	return func(v *Video) {
		v.Tags = tags
	}
}

func WithLanguage(language string) Option {
	return func(v *Video) {
		v.Language = language
	}
}

func WithLocale(locale string) Option {
	return func(v *Video) {
		v.Locale = locale
	}
}

func WithLicense(license string) Option {
	return func(v *Video) {
		v.License = license
	}
}

func WithThumbnail(thumbnail string) Option {
	return func(v *Video) {
		v.Thumbnail = thumbnail
	}
}

func WithRating(rating string) Option {
	return func(v *Video) {
		v.Rating = rating
	}
}

func WithChart(chart string) Option {
	return func(v *Video) {
		v.Chart = chart
	}
}

func WithForKids(forKids *bool) Option {
	return func(v *Video) {
		if forKids != nil {
			v.ForKids = forKids
		}
	}
}

func WithEmbeddable(embeddable *bool) Option {
	return func(v *Video) {
		if embeddable != nil {
			v.Embeddable = embeddable
		}
	}
}

func WithCategory(categoryId string) Option {
	return func(v *Video) {
		v.CategoryId = categoryId
	}
}

func WithPrivacy(privacy string) Option {
	return func(v *Video) {
		v.Privacy = privacy
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(v *Video) {
		v.PlaylistId = playlistId
	}
}

func WithPublicStatsViewable(publicStatsViewable *bool) Option {
	return func(v *Video) {
		if publicStatsViewable != nil {
			v.PublicStatsViewable = publicStatsViewable
		}
	}
}

func WithPublishAt(publishAt string) Option {
	return func(v *Video) {
		v.PublishAt = publishAt
	}
}

func WithRegionCode(regionCode string) Option {
	return func(v *Video) {
		v.RegionCode = regionCode
	}
}

func WithStabilize(stabilize *bool) Option {
	return func(v *Video) {
		if stabilize != nil {
			v.Stabilize = stabilize
		}
	}
}

func WithMaxHeight(maxHeight int64) Option {
	return func(v *Video) {
		v.MaxHeight = maxHeight
	}
}

func WithMaxWidth(maxWidth int64) Option {
	return func(v *Video) {
		v.MaxWidth = maxWidth
	}
}

func WithNotifySubscribers(notifySubscribers *bool) Option {
	return func(v *Video) {
		if notifySubscribers != nil {
			v.NotifySubscribers = notifySubscribers
		}
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(v *Video) {
		v.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithComments(comments string) Option {
	return func(v *Video) {
		v.Comments = comments
	}
}

func WithReasonId(reasonId string) Option {
	return func(v *Video) {
		v.ReasonId = reasonId
	}
}

func WithSecondaryReasonId(secondaryReasonId string) Option {
	return func(v *Video) {
		v.SecondaryReasonId = secondaryReasonId
	}
}

var (
	WithParts      = common.WithParts[*Video]
	WithOutput     = common.WithOutput[*Video]
	WithService    = common.WithService[*Video]
	WithIds        = common.WithIds[*Video]
	WithMaxResults = common.WithMaxResults[*Video]
	WithHl         = common.WithHl[*Video]
	WithChannelId  = common.WithChannelId[*Video]

	WithOnBehalfOfContentOwner = common.WithOnBehalfOfContentOwner[*Video]
)
