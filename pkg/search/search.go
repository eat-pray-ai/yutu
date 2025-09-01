package search

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	service      *youtube.Service
	errGetSearch = errors.New("failed to get search")
)

type search struct {
	ChannelId                 string   `yaml:"channel_id" json:"channel_id"`
	ChannelType               string   `yaml:"channel_type" json:"channel_type"`
	EventType                 string   `yaml:"event_type" json:"event_type"`
	ForContentOwner           *bool    `yaml:"for_content_owner" json:"for_content_owner"`
	ForDeveloper              *bool    `yaml:"for_developer" json:"for_developer"`
	ForMine                   *bool    `yaml:"for_mine" json:"for_mine"`
	Location                  string   `yaml:"location" json:"location"`
	LocationRadius            string   `yaml:"location_radius" json:"location_radius"`
	MaxResults                int64    `yaml:"max_results" json:"max_results"`
	OnBehalfOfContentOwner    string   `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
	Order                     string   `yaml:"order" json:"order"`
	PublishedAfter            string   `yaml:"published_after" json:"published_after"`
	PublishedBefore           string   `yaml:"published_before" json:"published_before"`
	Q                         string   `yaml:"q" json:"q"`
	RegionCode                string   `yaml:"region_code" json:"region_code"`
	RelevanceLanguage         string   `yaml:"relevance_language" json:"relevance_language"`
	SafeSearch                string   `yaml:"safe_search" json:"safe_search"`
	TopicId                   string   `yaml:"topic_id" json:"topic_id"`
	Types                     []string `yaml:"types" json:"types"`
	VideoCaption              string   `yaml:"video_caption" json:"video_caption"`
	VideoCategoryId           string   `yaml:"video_category_id" json:"video_category_id"`
	VideoDefinition           string   `yaml:"video_definition" json:"video_definition"`
	VideoDimension            string   `yaml:"video_dimension" json:"video_dimension"`
	VideoDuration             string   `yaml:"video_duration" json:"video_duration"`
	VideoEmbeddable           string   `yaml:"video_embeddable" json:"video_embeddable"`
	VideoLicense              string   `yaml:"video_license" json:"video_license"`
	VideoPaidProductPlacement string   `yaml:"video_paid_product_placement" json:"video_paid_product_placement"`
	VideoSyndicated           string   `yaml:"video_syndicated" json:"video_syndicated"`
	VideoType                 string   `yaml:"video_type" json:"video_type"`
}

type Search interface {
	Get([]string) ([]*youtube.SearchResult, error)
	List([]string, string, string, io.Writer) error
}

type Option func(*search)

func NewSearch(opts ...Option) Search {
	s := &search{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *search) Get(parts []string) ([]*youtube.SearchResult, error) {
	call := service.Search.List(parts)
	if s.ChannelId != "" {
		call.ChannelId(s.ChannelId)
	}

	if s.ChannelType != "" {
		call.ChannelType(s.ChannelType)
	}

	if s.EventType != "" {
		call.EventType(s.EventType)
	}

	if s.ForContentOwner != nil {
		call.ForContentOwner(*s.ForContentOwner)
	}

	if s.ForDeveloper != nil {
		call.ForDeveloper(*s.ForDeveloper)
	}

	if s.ForMine != nil {
		call.ForMine(*s.ForMine)
	}

	if s.Location != "" {
		call.Location(s.Location)
	}

	if s.LocationRadius != "" {
		call.LocationRadius(s.LocationRadius)
	}

	call.MaxResults(s.MaxResults)

	if s.OnBehalfOfContentOwner != "" {
		call.OnBehalfOfContentOwner(s.OnBehalfOfContentOwner)
	}

	if s.Order != "" {
		call.Order(s.Order)
	}

	if s.PublishedAfter != "" {
		call.PublishedAfter(s.PublishedAfter)
	}

	if s.PublishedBefore != "" {
		call.PublishedBefore(s.PublishedBefore)
	}

	if s.Q != "" {
		call.Q(s.Q)
	}

	if s.RegionCode != "" {
		call.RegionCode(s.RegionCode)
	}

	if s.RelevanceLanguage != "" {
		call.RelevanceLanguage(s.RelevanceLanguage)
	}

	if s.SafeSearch != "" {
		call.SafeSearch(s.SafeSearch)
	}

	if s.TopicId != "" {
		call.TopicId(s.TopicId)
	}

	if len(s.Types) > 0 {
		call.Type(s.Types...)
	}

	if s.VideoCaption != "" {
		call.VideoCaption(s.VideoCaption)
	}

	if s.VideoCategoryId != "" {
		call.VideoCategoryId(s.VideoCategoryId)
	}

	if s.VideoDefinition != "" {
		call.VideoDefinition(s.VideoDefinition)
	}

	if s.VideoDimension != "" {
		call.VideoDimension(s.VideoDimension)
	}

	if s.VideoDuration != "" {
		call.VideoDuration(s.VideoDuration)
	}

	if s.VideoEmbeddable != "" {
		call.VideoEmbeddable(s.VideoEmbeddable)
	}

	if s.VideoLicense != "" {
		call.VideoLicense(s.VideoLicense)
	}

	if s.VideoPaidProductPlacement != "" {
		call.VideoPaidProductPlacement(s.VideoPaidProductPlacement)
	}

	if s.VideoSyndicated != "" {
		call.VideoSyndicated(s.VideoSyndicated)
	}

	if s.VideoType != "" {
		call.VideoType(s.VideoType)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetSearch, err)
	}

	return res.Items, nil
}

func (s *search) List(
	parts []string, output string, jpath string, writer io.Writer,
) error {
	results, err := s.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(results, jpath, writer)
	case "yaml":
		utils.PrintYAML(results, jpath, writer)
	case "table":
		tb := table.NewWriter()
		defer tb.Render()
		tb.SetOutputMirror(writer)
		tb.SetStyle(table.StyleLight)
		tb.SetAutoIndex(true)
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
	return nil
}

func WithChannelId(channelId string) Option {
	return func(s *search) {
		s.ChannelId = channelId
	}
}

func WithChannelType(channelType string) Option {
	return func(s *search) {
		s.ChannelType = channelType
	}
}

func WithEventType(eventType string) Option {
	return func(s *search) {
		s.EventType = eventType
	}
}

func WithForContentOwner(forContentOwner *bool) Option {
	return func(s *search) {
		if forContentOwner != nil {
			s.ForContentOwner = forContentOwner
		}
	}
}

func WithForDeveloper(forDeveloper *bool) Option {
	return func(s *search) {
		if forDeveloper != nil {
			s.ForDeveloper = forDeveloper
		}
	}
}

func WithForMine(forMine *bool) Option {
	return func(s *search) {
		if forMine != nil {
			s.ForMine = forMine
		}
	}
}

func WithLocation(location string) Option {
	return func(s *search) {
		s.Location = location
	}
}

func WithLocationRadius(locationRadius string) Option {
	return func(s *search) {
		s.LocationRadius = locationRadius
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *search) {
		if maxResults <= 0 {
			maxResults = 1
		}
		s.MaxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(s *search) {
		s.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOrder(order string) Option {
	return func(s *search) {
		s.Order = order
	}
}

func WithPublishedAfter(publishedAfter string) Option {
	return func(s *search) {
		s.PublishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore string) Option {
	return func(s *search) {
		s.PublishedBefore = publishedBefore
	}
}

func WithQ(q string) Option {
	return func(s *search) {
		s.Q = q
	}
}

func WithRegionCode(regionCode string) Option {
	return func(s *search) {
		s.RegionCode = regionCode
	}
}

func WithRelevanceLanguage(relevanceLanguage string) Option {
	return func(s *search) {
		s.RelevanceLanguage = relevanceLanguage
	}
}

func WithSafeSearch(safeSearch string) Option {
	return func(s *search) {
		s.SafeSearch = safeSearch
	}
}

func WithTopicId(topicId string) Option {
	return func(s *search) {
		s.TopicId = topicId
	}
}

func WithTypes(types []string) Option {
	return func(s *search) {
		s.Types = types
	}
}

func WithVideoCaption(videoCaption string) Option {
	return func(s *search) {
		s.VideoCaption = videoCaption
	}
}

func WithVideoCategoryId(videoCategoryId string) Option {
	return func(s *search) {
		s.VideoCategoryId = videoCategoryId
	}
}

func WithVideoDefinition(videoDefinition string) Option {
	return func(s *search) {
		s.VideoDefinition = videoDefinition
	}
}

func WithVideoDimension(videoDimension string) Option {
	return func(s *search) {
		s.VideoDimension = videoDimension
	}
}

func WithVideoDuration(videoDuration string) Option {
	return func(s *search) {
		s.VideoDuration = videoDuration
	}
}

func WithVideoEmbeddable(videoEmbeddable string) Option {
	return func(s *search) {
		s.VideoEmbeddable = videoEmbeddable
	}
}

func WithVideoLicense(videoLicense string) Option {
	return func(s *search) {
		s.VideoLicense = videoLicense
	}
}

func WithVideoPaidProductPlacement(videoPaidProductPlacement string) Option {
	return func(s *search) {
		s.VideoPaidProductPlacement = videoPaidProductPlacement
	}
}

func WithVideoSyndicated(videoSyndicated string) Option {
	return func(s *search) {
		s.VideoSyndicated = videoSyndicated
	}
}

func WithVideoType(videoType string) Option {
	return func(s *search) {
		s.VideoType = videoType
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *search) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
