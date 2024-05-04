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
	errGetSearch error = errors.New("failed to get search")
)

type search struct {
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           bool
	forDeveloper              bool
	forMine                   bool
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

type searchOption func(*search)

func NewSearch(opts ...searchOption) Search {
	s := &search{}
	service = auth.NewY2BService()

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *search) get(parts []string) []*youtube.SearchResult {
	call := service.Search.List(parts)
	call.ChannelId(s.channelId)
	call.ChannelType(s.channelType)
	call.EventType(s.eventType)
	call.ForContentOwner(s.forContentOwner)
	call.ForDeveloper(s.forDeveloper)
	call.ForMine(s.forMine)
	call.Location(s.location)
	call.LocationRadius(s.locationRadius)
	call.MaxResults(s.maxResults)
	call.OnBehalfOfContentOwner(s.onBehalfOfContentOwner)
	call.Order(s.order)
	call.PublishedAfter(s.publishedAfter)
	call.PublishedBefore(s.publishedBefore)
	call.Q(s.q)
	call.RegionCode(s.regionCode)
	call.RelevanceLanguage(s.relevanceLanguage)
	call.SafeSearch(s.safeSearch)
	call.TopicId(s.topicId)
	call.Type(s.types)
	call.VideoCaption(s.videoCaption)
	call.VideoCategoryId(s.videoCategoryId)
	call.VideoDefinition(s.videoDefinition)
	call.VideoDimension(s.videoDimension)
	call.VideoDuration(s.videoDuration)
	call.VideoEmbeddable(s.videoEmbeddable)
	call.VideoLicense(s.videoLicense)
	call.VideoPaidProductPlacement(s.videoPaidProductPlacement)
	call.VideoSyndicated(s.videoSyndicated)
	call.VideoType(s.videoType)

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

func withSearchChannelId(channelId string) searchOption {
	return func(s *search) {
		s.channelId = channelId
	}
}

func withSearchChannelType(channelType string) searchOption {
	return func(s *search) {
		s.channelType = channelType
	}
}

func withSearchEventType(eventType string) searchOption {
	return func(s *search) {
		s.eventType = eventType
	}
}

func withSearchForContentOwner(forContentOwner bool) searchOption {
	return func(s *search) {
		s.forContentOwner = forContentOwner
	}
}

func withSearchForDeveloper(forDeveloper bool) searchOption {
	return func(s *search) {
		s.forDeveloper = forDeveloper
	}
}

func withSearchForMine(forMine bool) searchOption {
	return func(s *search) {
		s.forMine = forMine
	}
}

func withSearchLocation(location string) searchOption {
	return func(s *search) {
		s.location = location
	}
}

func withSearchLocationRadius(locationRadius string) searchOption {
	return func(s *search) {
		s.locationRadius = locationRadius
	}
}

func withSearchMaxResults(maxResults int64) searchOption {
	return func(s *search) {
		s.maxResults = maxResults
	}
}

func withSearchOnBehalfOfContentOwner(onBehalfOfContentOwner string) searchOption {
	return func(s *search) {
		s.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func withSearchOrder(order string) searchOption {
	return func(s *search) {
		s.order = order
	}
}

func withSearchPublishedAfter(publishedAfter string) searchOption {
	return func(s *search) {
		s.publishedAfter = publishedAfter
	}
}

func withSearchPublishedBefore(publishedBefore string) searchOption {
	return func(s *search) {
		s.publishedBefore = publishedBefore
	}
}

func withSearchQ(q string) searchOption {
	return func(s *search) {
		s.q = q
	}
}

func withSearchRegionCode(regionCode string) searchOption {
	return func(s *search) {
		s.regionCode = regionCode
	}
}

func withSearchRelevanceLanguage(relevanceLanguage string) searchOption {
	return func(s *search) {
		s.relevanceLanguage = relevanceLanguage
	}
}

func withSearchSafeSearch(safeSearch string) searchOption {
	return func(s *search) {
		s.safeSearch = safeSearch
	}
}

func withSearchTopicId(topicId string) searchOption {
	return func(s *search) {
		s.topicId = topicId
	}
}

func withSearchTypes(types string) searchOption {
	return func(s *search) {
		s.types = types
	}
}

func withSearchVideoCaption(videoCaption string) searchOption {
	return func(s *search) {
		s.videoCaption = videoCaption
	}
}

func withSearchVideoCategoryId(videoCategoryId string) searchOption {
	return func(s *search) {
		s.videoCategoryId = videoCategoryId
	}
}

func withSearchVideoDefinition(videoDefinition string) searchOption {
	return func(s *search) {
		s.videoDefinition = videoDefinition
	}
}

func withSearchVideoDimension(videoDimension string) searchOption {
	return func(s *search) {
		s.videoDimension = videoDimension
	}
}

func withSearchVideoDuration(videoDuration string) searchOption {
	return func(s *search) {
		s.videoDuration = videoDuration
	}
}

func withSearchVideoEmbeddable(videoEmbeddable string) searchOption {
	return func(s *search) {
		s.videoEmbeddable = videoEmbeddable
	}
}

func withSearchVideoLicense(videoLicense string) searchOption {
	return func(s *search) {
		s.videoLicense = videoLicense
	}
}

func withSearchVideoPaidProductPlacement(videoPaidProductPlacement string) searchOption {
	return func(s *search) {
		s.videoPaidProductPlacement = videoPaidProductPlacement
	}
}

func withSearchVideoSyndicated(videoSyndicated string) searchOption {
	return func(s *search) {
		s.videoSyndicated = videoSyndicated
	}
}

func withSearchVideoType(videoType string) searchOption {
	return func(s *search) {
		s.videoType = videoType
	}
}
