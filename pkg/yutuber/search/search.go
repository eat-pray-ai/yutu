package search

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service      *youtube.Service
	errGetSearch = errors.New("failed to get search")
)

type search struct {
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           *bool
	forDeveloper              *bool
	forMine                   *bool
	location                  string
	locationRadius            string
	maxResults                int64
	onBehalfOfContentOwner    string
	order                     string
	publishedAfter            string
	publishedBefore           string
	q                         string
	regionCode                string
	relevanceLanguage         string
	safeSearch                string
	topicId                   string
	types                     string
	videoCaption              string
	videoCategoryId           string
	videoDefinition           string
	videoDimension            string
	videoDuration             string
	videoEmbeddable           string
	videoLicense              string
	videoPaidProductPlacement string
	videoSyndicated           string
	videoType                 string
}

type Search interface {
	get([]string) []*youtube.SearchResult
	List([]string, string)
}

type Option func(*search)

func NewSearch(opts ...Option) Search {
	s := &search{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *search) get(parts []string) []*youtube.SearchResult {
	call := service.Search.List(parts)
	if s.channelId != "" {
		call.ChannelId(s.channelId)
	}

	if s.channelType != "" {
		call.ChannelType(s.channelType)
	}

	if s.eventType != "" {
		call.EventType(s.eventType)
	}

	if s.forContentOwner != nil {
		call.ForContentOwner(*s.forContentOwner)
	}

	if s.forDeveloper != nil {
		call.ForDeveloper(*s.forDeveloper)
	}

	if s.forMine != nil {
		call.ForMine(*s.forMine)
	}

	if s.location != "" {
		call.Location(s.location)
	}

	if s.locationRadius != "" {
		call.LocationRadius(s.locationRadius)
	}

	call.MaxResults(s.maxResults)

	if s.onBehalfOfContentOwner != "" {
		call.OnBehalfOfContentOwner(s.onBehalfOfContentOwner)
	}

	if s.order != "" {
		call.Order(s.order)
	}

	if s.publishedAfter != "" {
		call.PublishedAfter(s.publishedAfter)
	}

	if s.publishedBefore != "" {
		call.PublishedBefore(s.publishedBefore)
	}

	if s.q != "" {
		call.Q(s.q)
	}

	if s.regionCode != "" {
		call.RegionCode(s.regionCode)
	}

	if s.relevanceLanguage != "" {
		call.RelevanceLanguage(s.relevanceLanguage)
	}

	if s.safeSearch != "" {
		call.SafeSearch(s.safeSearch)
	}

	if s.topicId != "" {
		call.TopicId(s.topicId)
	}

	if s.types != "" {
		call.Type(s.types)
	}

	if s.videoCaption != "" {
		call.VideoCaption(s.videoCaption)
	}

	if s.videoCategoryId != "" {
		call.VideoCategoryId(s.videoCategoryId)
	}

	if s.videoDefinition != "" {
		call.VideoDefinition(s.videoDefinition)
	}

	if s.videoDimension != "" {
		call.VideoDimension(s.videoDimension)
	}

	if s.videoDuration != "" {
		call.VideoDuration(s.videoDuration)
	}

	if s.videoEmbeddable != "" {
		call.VideoEmbeddable(s.videoEmbeddable)
	}

	if s.videoLicense != "" {
		call.VideoLicense(s.videoLicense)
	}

	if s.videoPaidProductPlacement != "" {
		call.VideoPaidProductPlacement(s.videoPaidProductPlacement)
	}

	if s.videoSyndicated != "" {
		call.VideoSyndicated(s.videoSyndicated)
	}

	if s.videoType != "" {
		call.VideoType(s.videoType)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetSearch, err), s.q)
	}

	return res.Items
}

func (s *search) List(parts []string, output string) {
	results := s.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(results)
	case "yaml":
		utils.PrintYAML(results)
	default:
		fmt.Println("Kind\tTitle")
		for _, result := range results {
			fmt.Printf("%v\t%v\n", result.Id.Kind, result.Snippet.Title)
		}
	}
}

func WithChannelId(channelId string) Option {
	return func(s *search) {
		s.channelId = channelId
	}
}

func WithChannelType(channelType string) Option {
	return func(s *search) {
		s.channelType = channelType
	}
}

func WithEventType(eventType string) Option {
	return func(s *search) {
		s.eventType = eventType
	}
}

func WithForContentOwner(forContentOwner bool, changed bool) Option {
	return func(s *search) {
		if changed {
			s.forContentOwner = &forContentOwner
		}
	}
}

func WithForDeveloper(forDeveloper bool, changed bool) Option {
	return func(s *search) {
		if changed {
			s.forDeveloper = &forDeveloper
		}
	}
}

func WithForMine(forMine bool, changed bool) Option {
	return func(s *search) {
		if changed {
			s.forMine = &forMine
		}
	}
}

func WithLocation(location string) Option {
	return func(s *search) {
		s.location = location
	}
}

func WithLocationRadius(locationRadius string) Option {
	return func(s *search) {
		s.locationRadius = locationRadius
	}
}

func WithMaxResults(maxResults int64) Option {
	return func(s *search) {
		s.maxResults = maxResults
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(s *search) {
		s.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithOrder(order string) Option {
	return func(s *search) {
		s.order = order
	}
}

func WithPublishedAfter(publishedAfter string) Option {
	return func(s *search) {
		s.publishedAfter = publishedAfter
	}
}

func WithPublishedBefore(publishedBefore string) Option {
	return func(s *search) {
		s.publishedBefore = publishedBefore
	}
}

func WithQ(q string) Option {
	return func(s *search) {
		s.q = q
	}
}

func WithRegionCode(regionCode string) Option {
	return func(s *search) {
		s.regionCode = regionCode
	}
}

func WithRelevanceLanguage(relevanceLanguage string) Option {
	return func(s *search) {
		s.relevanceLanguage = relevanceLanguage
	}
}

func WithSafeSearch(safeSearch string) Option {
	return func(s *search) {
		s.safeSearch = safeSearch
	}
}

func WithTopicId(topicId string) Option {
	return func(s *search) {
		s.topicId = topicId
	}
}

func WithTypes(types string) Option {
	return func(s *search) {
		s.types = types
	}
}

func WithVideoCaption(videoCaption string) Option {
	return func(s *search) {
		s.videoCaption = videoCaption
	}
}

func WithVideoCategoryId(videoCategoryId string) Option {
	return func(s *search) {
		s.videoCategoryId = videoCategoryId
	}
}

func WithVideoDefinition(videoDefinition string) Option {
	return func(s *search) {
		s.videoDefinition = videoDefinition
	}
}

func WithVideoDimension(videoDimension string) Option {
	return func(s *search) {
		s.videoDimension = videoDimension
	}
}

func WithVideoDuration(videoDuration string) Option {
	return func(s *search) {
		s.videoDuration = videoDuration
	}
}

func WithVideoEmbeddable(videoEmbeddable string) Option {
	return func(s *search) {
		s.videoEmbeddable = videoEmbeddable
	}
}

func WithVideoLicense(videoLicense string) Option {
	return func(s *search) {
		s.videoLicense = videoLicense
	}
}

func WithVideoPaidProductPlacement(videoPaidProductPlacement string) Option {
	return func(s *search) {
		s.videoPaidProductPlacement = videoPaidProductPlacement
	}
}

func WithVideoSyndicated(videoSyndicated string) Option {
	return func(s *search) {
		s.videoSyndicated = videoSyndicated
	}
}

func WithVideoType(videoType string) Option {
	return func(s *search) {
		s.videoType = videoType
	}
}

func WithService() Option {
	return func(s *search) {
		service = auth.NewY2BService()
	}
}
