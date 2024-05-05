package yutuber

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	errGetSearch = errors.New("failed to get search")
)

type search struct {
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           string
	forDeveloper              string
	forMine                   string
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

type SearchOption func(*search)

func NewSearch(opts ...SearchOption) Search {
	s := &search{}
	service = auth.NewY2BService()

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

	if s.forContentOwner == "true" {
		call.ForContentOwner(true)
	} else if s.forContentOwner == "false" {
		call.ForContentOwner(false)
	}

	if s.forDeveloper == "true" {
		call.ForDeveloper(true)
	} else if s.forDeveloper == "false" {
		call.ForDeveloper(false)
	}

	if s.forMine == "true" {
		call.ForMine(true)
	} else if s.forMine == "false" {
		call.ForMine(false)
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

	resp, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetSearch, err), s.q)
	}

	return resp.Items
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

func WithSearchChannelId(channelId string) SearchOption {
	return func(s *search) {
		s.channelId = channelId
	}
}

func WithSearchChannelType(channelType string) SearchOption {
	return func(s *search) {
		s.channelType = channelType
	}
}

func WithSearchEventType(eventType string) SearchOption {
	return func(s *search) {
		s.eventType = eventType
	}
}

func WithSearchForContentOwner(forContentOwner string) SearchOption {
	return func(s *search) {
		s.forContentOwner = forContentOwner
	}
}

func WithSearchForDeveloper(forDeveloper string) SearchOption {
	return func(s *search) {
		s.forDeveloper = forDeveloper
	}
}

func WithSearchForMine(forMine string) SearchOption {
	return func(s *search) {
		s.forMine = forMine
	}
}

func WithSearchLocation(location string) SearchOption {
	return func(s *search) {
		s.location = location
	}
}

func WithSearchLocationRadius(locationRadius string) SearchOption {
	return func(s *search) {
		s.locationRadius = locationRadius
	}
}

func WithSearchMaxResults(maxResults int64) SearchOption {
	return func(s *search) {
		s.maxResults = maxResults
	}
}

func WithSearchOnBehalfOfContentOwner(onBehalfOfContentOwner string) SearchOption {
	return func(s *search) {
		s.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithSearchOrder(order string) SearchOption {
	return func(s *search) {
		s.order = order
	}
}

func WithSearchPublishedAfter(publishedAfter string) SearchOption {
	return func(s *search) {
		s.publishedAfter = publishedAfter
	}
}

func WithSearchPublishedBefore(publishedBefore string) SearchOption {
	return func(s *search) {
		s.publishedBefore = publishedBefore
	}
}

func WithSearchQ(q string) SearchOption {
	return func(s *search) {
		s.q = q
	}
}

func WithSearchRegionCode(regionCode string) SearchOption {
	return func(s *search) {
		s.regionCode = regionCode
	}
}

func WithSearchRelevanceLanguage(relevanceLanguage string) SearchOption {
	return func(s *search) {
		s.relevanceLanguage = relevanceLanguage
	}
}

func WithSearchSafeSearch(safeSearch string) SearchOption {
	return func(s *search) {
		s.safeSearch = safeSearch
	}
}

func WithSearchTopicId(topicId string) SearchOption {
	return func(s *search) {
		s.topicId = topicId
	}
}

func WithSearchTypes(types string) SearchOption {
	return func(s *search) {
		s.types = types
	}
}

func WithSearchVideoCaption(videoCaption string) SearchOption {
	return func(s *search) {
		s.videoCaption = videoCaption
	}
}

func WithSearchVideoCategoryId(videoCategoryId string) SearchOption {
	return func(s *search) {
		s.videoCategoryId = videoCategoryId
	}
}

func WithSearchVideoDefinition(videoDefinition string) SearchOption {
	return func(s *search) {
		s.videoDefinition = videoDefinition
	}
}

func WithSearchVideoDimension(videoDimension string) SearchOption {
	return func(s *search) {
		s.videoDimension = videoDimension
	}
}

func WithSearchVideoDuration(videoDuration string) SearchOption {
	return func(s *search) {
		s.videoDuration = videoDuration
	}
}

func WithSearchVideoEmbeddable(videoEmbeddable string) SearchOption {
	return func(s *search) {
		s.videoEmbeddable = videoEmbeddable
	}
}

func WithSearchVideoLicense(videoLicense string) SearchOption {
	return func(s *search) {
		s.videoLicense = videoLicense
	}
}

func WithSearchVideoPaidProductPlacement(videoPaidProductPlacement string) SearchOption {
	return func(s *search) {
		s.videoPaidProductPlacement = videoPaidProductPlacement
	}
}

func WithSearchVideoSyndicated(videoSyndicated string) SearchOption {
	return func(s *search) {
		s.videoSyndicated = videoSyndicated
	}
}

func WithSearchVideoType(videoType string) SearchOption {
	return func(s *search) {
		s.videoType = videoType
	}
}
