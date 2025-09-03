package video

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

var (
	service        *youtube.Service
	errGetVideo    = errors.New("failed to get video")
	errInsertVideo = errors.New("failed to insert video")
	errUpdateVideo = errors.New("failed to update video")
	errRating      = errors.New("failed to rate video")
	errGetRating   = errors.New("failed to get rating")
	errDeleteVideo = errors.New("failed to delete video")
	errReportAbuse = errors.New("failed to report abuse")
)

type video struct {
	IDs               []string `yaml:"ids" json:"ids"`
	AutoLevels        *bool    `yaml:"auto_levels" json:"auto_levels"`
	File              string   `yaml:"file" json:"file"`
	Title             string   `yaml:"title" json:"title"`
	Description       string   `yaml:"description" json:"description"`
	Hl                string   `yaml:"hl" json:"hl"`
	Tags              []string `yaml:"tags" json:"tags"`
	Language          string   `yaml:"language" json:"language"`
	Locale            string   `yaml:"locale" json:"locale"`
	License           string   `yaml:"license" json:"license"`
	Thumbnail         string   `yaml:"thumbnail" json:"thumbnail"`
	Rating            string   `yaml:"rating" json:"rating"`
	Chart             string   `yaml:"chart" json:"chart"`
	ChannelId         string   `yaml:"channel_id" json:"channel_id"`
	Comments          string   `yaml:"comments" json:"comments"`
	PlaylistId        string   `yaml:"playlist_id" json:"playlist_id"`
	CategoryId        string   `yaml:"category_id" json:"category_id"`
	Privacy           string   `yaml:"privacy" json:"privacy"`
	ForKids           *bool    `yaml:"for_kids" json:"for_kids"`
	Embeddable        *bool    `yaml:"embeddable" json:"embeddable"`
	PublishAt         string   `yaml:"publish_at" json:"publish_at"`
	RegionCode        string   `yaml:"region_code" json:"region_code"`
	ReasonId          string   `yaml:"reason_id" json:"reason_id"`
	SecondaryReasonId string   `yaml:"secondary_reason_id" json:"secondary_reason_id"`
	Stabilize         *bool    `yaml:"stabilize" json:"stabilize"`
	MaxHeight         int64    `yaml:"max_height" json:"max_height"`
	MaxWidth          int64    `yaml:"max_width" json:"max_width"`
	MaxResults        int64    `yaml:"max_results" json:"max_results"`

	NotifySubscribers             *bool  `yaml:"notify_subscribers" json:"notify_subscribers"`
	PublicStatsViewable           *bool  `yaml:"public_stats_viewable" json:"public_stats_viewable"`
	OnBehalfOfContentOwner        string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	OnBehalfOfContentOwnerChannel string `yaml:"on_behalf_of_content_owner_channel" json:"on_behalf_of_content_owner_channel"`
}

type Video interface {
	List([]string, string, string, io.Writer) error
	Insert(string, string, io.Writer) error
	Update(string, string, io.Writer) error
	Rate(io.Writer) error
	GetRating(string, string, io.Writer) error
	Delete(io.Writer) error
	ReportAbuse(io.Writer) error
	Get([]string) ([]*youtube.Video, error)
}

type Option func(*video)

func NewVideo(opts ...Option) Video {
	v := &video{}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *video) Get(parts []string) ([]*youtube.Video, error) {
	call := service.Videos.List(parts)
	if len(v.IDs) > 0 {
		call = call.Id(v.IDs...)
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
	call = call.MaxResults(v.MaxResults)

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideo, err)
	}

	return res.Items, nil
}

func (v *video) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	videos, err := v.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(videos, jpath, writer)
	case "yaml":
		utils.PrintYAML(videos, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Title", "Channel ID", "Views"})
		for _, video := range videos {
			tb.AppendRow(
				table.Row{
					video.Id, video.Snippet.Title,
					video.Snippet.ChannelId,
					video.Statistics.ViewCount,
				},
			)
		}
	}
	return nil
}

func (v *video) Insert(output string, jpath string, writer io.Writer) error {
	file, err := pkg.Root.Open(v.File)
	if err != nil {
		return errors.Join(errInsertVideo, err)
	}
	defer file.Close()

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

	call := service.Videos.Insert([]string{"snippet,status"}, video)

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
			thumbnail.WithService(service),
		)
		_ = t.Set("silent", "", nil)
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
			playlistItem.WithService(service),
		)

		_ = pi.Insert("silent", "", writer)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Video inserted: %s\n", res.Id)
	}
	return nil
}

func (v *video) Update(output string, jpath string, writer io.Writer) error {
	videos, err := v.Get([]string{"id", "snippet", "status"})
	if err != nil {
		return errors.Join(errUpdateVideo, err)
	}
	if len(videos) == 0 {
		return errGetVideo
	}

	video := videos[0]
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
		video.Snippet.DefaultAudioLanguage = v.Language
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

	call := service.Videos.Update([]string{"snippet,status"}, video)
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
			thumbnail.WithService(service),
		)
		_ = t.Set("silent", "", nil)
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
			playlistItem.WithService(service),
		)

		_ = pi.Insert("silent", "", writer)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Video updated: %s\n", res.Id)
	}
	return nil
}

func (v *video) Rate(writer io.Writer) error {
	for _, id := range v.IDs {
		call := service.Videos.Rate(id, v.Rating)
		err := call.Do()
		if err != nil {
			return errors.Join(errRating, err)
		}
		_, _ = fmt.Fprintf(writer, "Video %s rated %s\n", id, v.Rating)
	}
	return nil
}

func (v *video) GetRating(output string, jpath string, writer io.Writer) error {
	call := service.Videos.GetRating(v.IDs)
	if v.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(v.OnBehalfOfContentOwnerChannel)
	}
	res, err := call.Do()
	if err != nil {
		return errors.Join(errGetRating, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res.Items, jpath, writer)
	case "yaml":
		utils.PrintYAML(res.Items, jpath, writer)
	default:
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
		tb.AppendHeader(table.Row{"ID", "Rating"})
		for _, item := range res.Items {
			tb.AppendRow(table.Row{item.VideoId, item.Rating})
		}
	}
	return nil
}

func (v *video) Delete(writer io.Writer) error {
	for _, id := range v.IDs {
		call := service.Videos.Delete(id)
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

func (v *video) ReportAbuse(writer io.Writer) error {
	for _, id := range v.IDs {
		videoAbuseReport := &youtube.VideoAbuseReport{
			Comments:          v.Comments,
			Language:          v.Language,
			ReasonId:          v.ReasonId,
			SecondaryReasonId: v.SecondaryReasonId,
			VideoId:           id,
		}

		call := service.Videos.ReportAbuse(videoAbuseReport)
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

func WithIDs(ids []string) Option {
	return func(v *video) {
		v.IDs = ids
	}
}

func WithAutoLevels(autoLevels *bool) Option {
	return func(v *video) {
		if autoLevels != nil {
			v.AutoLevels = autoLevels
		}
	}
}

func WithFile(file string) Option {
	return func(v *video) {
		v.File = file
	}
}

func WithTitle(title string) Option {
	return func(v *video) {
		v.Title = title
	}
}

func WithDescription(description string) Option {
	return func(v *video) {
		v.Description = description
	}
}

func WithHl(hl string) Option {
	return func(v *video) {
		v.Hl = hl
	}
}

func WithTags(tags []string) Option {
	return func(v *video) {
		v.Tags = tags
	}
}

func WithLanguage(language string) Option {
	return func(v *video) {
		v.Language = language
	}
}

func WithLocale(locale string) Option {
	return func(v *video) {
		v.Locale = locale
	}
}

func WithLicense(license string) Option {
	return func(v *video) {
		v.License = license
	}
}

func WithThumbnail(thumbnail string) Option {
	return func(v *video) {
		v.Thumbnail = thumbnail
	}
}

func WithRating(rating string) Option {
	return func(v *video) {
		v.Rating = rating
	}
}

func WithChart(chart string) Option {
	return func(v *video) {
		v.Chart = chart
	}
}

func WithForKids(forKids *bool) Option {
	return func(v *video) {
		if forKids != nil {
			v.ForKids = forKids
		}
	}
}

func WithEmbeddable(embeddable *bool) Option {
	return func(v *video) {
		if embeddable != nil {
			v.Embeddable = embeddable
		}
	}
}

func WithCategory(categoryId string) Option {
	return func(v *video) {
		v.CategoryId = categoryId
	}
}

func WithPrivacy(privacy string) Option {
	return func(v *video) {
		v.Privacy = privacy
	}
}

func WithChannelId(channelId string) Option {
	return func(v *video) {
		v.ChannelId = channelId
	}
}

func WithPlaylistId(playlistId string) Option {
	return func(v *video) {
		v.PlaylistId = playlistId
	}
}

func WithPublicStatsViewable(publicStatsViewable *bool) Option {
	return func(v *video) {
		if publicStatsViewable != nil {
			v.PublicStatsViewable = publicStatsViewable
		}
	}
}

func WithPublishAt(publishAt string) Option {
	return func(v *video) {
		v.PublishAt = publishAt
	}
}

func WithRegionCode(regionCode string) Option {
	return func(v *video) {
		v.RegionCode = regionCode
	}
}

func WithStabilize(stabilize *bool) Option {
	return func(v *video) {
		if stabilize != nil {
			v.Stabilize = stabilize
		}
	}
}

func WithMaxHeight(maxHeight int64) Option {
	return func(v *video) {
		v.MaxHeight = maxHeight
	}
}

func WithMaxWidth(maxWidth int64) Option {
	return func(v *video) {
		v.MaxWidth = maxWidth
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(v *video) {
		if maxResults <= 0 {
			maxResults = 1
		}
		v.MaxResults = maxResults
	}
}

func WithNotifySubscribers(notifySubscribers *bool) Option {
	return func(v *video) {
		if notifySubscribers != nil {
			v.NotifySubscribers = notifySubscribers
		}
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(v *video) {
		v.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel string) Option {
	return func(v *video) {
		v.OnBehalfOfContentOwnerChannel = onBehalfOfContentOwnerChannel
	}
}

func WithComments(comments string) Option {
	return func(v *video) {
		v.Comments = comments
	}
}

func WithReasonId(reasonId string) Option {
	return func(v *video) {
		v.ReasonId = reasonId
	}
}

func WithSecondaryReasonId(secondaryReasonId string) Option {
	return func(v *video) {
		v.SecondaryReasonId = secondaryReasonId
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *video) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}
