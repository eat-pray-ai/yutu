// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"errors"
	"io"
	"math"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetSearch = errors.New("failed to get search")
)

type Search struct {
	*common.Fields
	ChannelId                 string   `yaml:"channel_id" json:"channel_id,omitempty"`
	ChannelType               string   `yaml:"channel_type" json:"channel_type,omitempty"`
	EventType                 string   `yaml:"event_type" json:"event_type,omitempty"`
	ForContentOwner           *bool    `yaml:"for_content_owner" json:"for_content_owner,omitempty"`
	ForDeveloper              *bool    `yaml:"for_developer" json:"for_developer,omitempty"`
	ForMine                   *bool    `yaml:"for_mine" json:"for_mine,omitempty"`
	Location                  string   `yaml:"location" json:"location,omitempty"`
	LocationRadius            string   `yaml:"location_radius" json:"location_radius,omitempty"`
	MaxResults                int64    `yaml:"max_results" json:"max_results,omitempty"`
	OnBehalfOfContentOwner    string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner,omitempty"`
	Order                     string   `yaml:"order" json:"order,omitempty"`
	PublishedAfter            string   `yaml:"published_after" json:"published_after,omitempty"`
	PublishedBefore           string   `yaml:"published_before" json:"published_before,omitempty"`
	Q                         string   `yaml:"q" json:"q,omitempty"`
	RegionCode                string   `yaml:"region_code" json:"region_code,omitempty"`
	RelevanceLanguage         string   `yaml:"relevance_language" json:"relevance_language,omitempty"`
	SafeSearch                string   `yaml:"safe_search" json:"safe_search,omitempty"`
	TopicId                   string   `yaml:"topic_id" json:"topic_id,omitempty"`
	Types                     []string `yaml:"types" json:"types,omitempty"`
	VideoCaption              string   `yaml:"video_caption" json:"video_caption,omitempty"`
	VideoCategoryId           string   `yaml:"video_category_id" json:"video_category_id,omitempty"`
	VideoDefinition           string   `yaml:"video_definition" json:"video_definition,omitempty"`
	VideoDimension            string   `yaml:"video_dimension" json:"video_dimension,omitempty"`
	VideoDuration             string   `yaml:"video_duration" json:"video_duration,omitempty"`
	VideoEmbeddable           string   `yaml:"video_embeddable" json:"video_embeddable,omitempty"`
	VideoLicense              string   `yaml:"video_license" json:"video_license,omitempty"`
	VideoPaidProductPlacement string   `yaml:"video_paid_product_placement" json:"video_paid_product_placement,omitempty"`
	VideoSyndicated           string   `yaml:"video_syndicated" json:"video_syndicated,omitempty"`
	VideoType                 string   `yaml:"video_type" json:"video_type,omitempty"`
}

type ISearch[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*Search)

func NewSearch(opts ...Option) ISearch[youtube.SearchResult] {
	s := &Search{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Search) Get() ([]*youtube.SearchResult, error) {
	s.EnsureService()
	call := s.Service.Search.List(s.Parts)
	if s.ChannelId != "" {
		call = call.ChannelId(s.ChannelId)
	}
	if s.ChannelType != "" {
		call = call.ChannelType(s.ChannelType)
	}
	if s.EventType != "" {
		call = call.EventType(s.EventType)
	}
	if s.ForContentOwner != nil {
		call = call.ForContentOwner(*s.ForContentOwner)
	}
	if s.ForDeveloper != nil {
		call = call.ForDeveloper(*s.ForDeveloper)
	}
	if s.ForMine != nil {
		call = call.ForMine(*s.ForMine)
	}
	if s.Location != "" {
		call = call.Location(s.Location)
	}
	if s.LocationRadius != "" {
		call = call.LocationRadius(s.LocationRadius)
	}
	if s.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}
	if s.Order != "" {
		call = call.Order(s.Order)
	}
	if s.PublishedAfter != "" {
		call = call.PublishedAfter(s.PublishedAfter)
	}
	if s.PublishedBefore != "" {
		call = call.PublishedBefore(s.PublishedBefore)
	}
	if s.Q != "" {
		call = call.Q(s.Q)
	}
	if s.RegionCode != "" {
		call = call.RegionCode(s.RegionCode)
	}
	if s.RelevanceLanguage != "" {
		call = call.RelevanceLanguage(s.RelevanceLanguage)
	}
	if s.SafeSearch != "" {
		call = call.SafeSearch(s.SafeSearch)
	}
	if s.TopicId != "" {
		call = call.TopicId(s.TopicId)
	}
	if len(s.Types) > 0 {
		call = call.Type(s.Types...)
	}
	if s.VideoCaption != "" {
		call = call.VideoCaption(s.VideoCaption)
	}
	if s.VideoCategoryId != "" {
		call = call.VideoCategoryId(s.VideoCategoryId)
	}
	if s.VideoDefinition != "" {
		call = call.VideoDefinition(s.VideoDefinition)
	}
	if s.VideoDimension != "" {
		call = call.VideoDimension(s.VideoDimension)
	}
	if s.VideoDuration != "" {
		call = call.VideoDuration(s.VideoDuration)
	}
	if s.VideoEmbeddable != "" {
		call = call.VideoEmbeddable(s.VideoEmbeddable)
	}
	if s.VideoLicense != "" {
		call = call.VideoLicense(s.VideoLicense)
	}
	if s.VideoPaidProductPlacement != "" {
		call = call.VideoPaidProductPlacement(s.VideoPaidProductPlacement)
	}
	if s.VideoSyndicated != "" {
		call = call.VideoSyndicated(s.VideoSyndicated)
	}
	if s.VideoType != "" {
		call = call.VideoType(s.VideoType)
	}

	var items []*youtube.SearchResult
	pageToken := ""
	for s.MaxResults > 0 {
		call = call.MaxResults(min(s.MaxResults, pkg.PerPage))
		s.MaxResults -= pkg.PerPage
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		res, err := call.Do()
		if err != nil {
			return items, errors.Join(errGetSearch, err)
		}

		items = append(items, res.Items...)
		pageToken = res.NextPageToken
		if pageToken == "" || len(res.Items) == 0 {
			break
		}
	}

	return items, nil
}

func (s *Search) List(writer io.Writer) error {
	results, err := s.Get()
	if err != nil && results == nil {
		return err
	}

	switch s.Output {
	case "json":
		utils.PrintJSON(results, s.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(results, s.Jsonpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(pkg.TableStyle)
		tb.AppendHeader(table.Row{"Kind", "Title", "Resource ID"})
		for _, result := range results {
			var resourceId string
			switch result.Id.Kind {
			case "youtube#video":
				resourceId = result.Id.VideoId
			case "youtube#channel":
				resourceId = result.Id.ChannelId
			case "youtube#playlist":
				resourceId = result.Id.PlaylistId
			}
			tb.AppendRow(
				table.Row{result.Id.Kind, result.Snippet.Title, resourceId},
			)
		}
	}
	return err
}

func WithChannelId(channelId string) Option {
	return func(s *Search) {
		s.ChannelId = channelId
	}
}

func WithChannelType(channelType string) Option {
	return func(s *Search) {
		s.ChannelType = channelType
	}
}

func WithEventType(eventType string) Option {
	return func(s *Search) {
		s.EventType = eventType
	}
}

func WithForContentOwner(forContentOwner *bool) Option {
	return func(s *Search) {
		if forContentOwner != nil {
			s.ForContentOwner = forContentOwner
		}
	}
}

func WithForDeveloper(forDeveloper *bool) Option {
	return func(s *Search) {
		if forDeveloper != nil {
			s.ForDeveloper = forDeveloper
		}
	}
}

func WithForMine(forMine *bool) Option {
	return func(s *Search) {
		if forMine != nil {
			s.ForMine = forMine
		}
	}
}

func WithLocation(location string) Option {
	return func(s *Search) {
		s.Location = location
	}
}

func WithLocationRadius(locationRadius string) Option {
	return func(s *Search) {
		s.LocationRadius = locationRadius
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *Search) {
		if maxResults < 0 {
			maxResults = 1
		} else if maxResults == 0 {
			maxResults = math.MaxInt64
		}
		s.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(s *Search) {
		s.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOrder(order string) Option {
	return func(s *Search) {
		s.Order = order
	}
}

func WithPublishedAfter(publishedAfter string) Option {
	return func(s *Search) {
		s.PublishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore string) Option {
	return func(s *Search) {
		s.PublishedBefore = publishedBefore
	}
}

func WithQ(q string) Option {
	return func(s *Search) {
		s.Q = q
	}
}

func WithRegionCode(regionCode string) Option {
	return func(s *Search) {
		s.RegionCode = regionCode
	}
}

func WithRelevanceLanguage(relevanceLanguage string) Option {
	return func(s *Search) {
		s.RelevanceLanguage = relevanceLanguage
	}
}

func WithSafeSearch(safeSearch string) Option {
	return func(s *Search) {
		s.SafeSearch = safeSearch
	}
}

func WithTopicId(topicId string) Option {
	return func(s *Search) {
		s.TopicId = topicId
	}
}

func WithTypes(types []string) Option {
	return func(s *Search) {
		s.Types = types
	}
}

func WithVideoCaption(videoCaption string) Option {
	return func(s *Search) {
		s.VideoCaption = videoCaption
	}
}

func WithVideoCategoryId(videoCategoryId string) Option {
	return func(s *Search) {
		s.VideoCategoryId = videoCategoryId
	}
}

func WithVideoDefinition(videoDefinition string) Option {
	return func(s *Search) {
		s.VideoDefinition = videoDefinition
	}
}

func WithVideoDimension(videoDimension string) Option {
	return func(s *Search) {
		s.VideoDimension = videoDimension
	}
}

func WithVideoDuration(videoDuration string) Option {
	return func(s *Search) {
		s.VideoDuration = videoDuration
	}
}

func WithVideoEmbeddable(videoEmbeddable string) Option {
	return func(s *Search) {
		s.VideoEmbeddable = videoEmbeddable
	}
}

func WithVideoLicense(videoLicense string) Option {
	return func(s *Search) {
		s.VideoLicense = videoLicense
	}
}

func WithVideoPaidProductPlacement(videoPaidProductPlacement string) Option {
	return func(s *Search) {
		s.VideoPaidProductPlacement = videoPaidProductPlacement
	}
}

func WithVideoSyndicated(videoSyndicated string) Option {
	return func(s *Search) {
		s.VideoSyndicated = videoSyndicated
	}
}

func WithVideoType(videoType string) Option {
	return func(s *Search) {
		s.VideoType = videoType
	}
}

var (
	WithParts    = common.WithParts[*Search]
	WithOutput   = common.WithOutput[*Search]
	WithJsonpath = common.WithJsonpath[*Search]
	WithService  = common.WithService[*Search]
)
