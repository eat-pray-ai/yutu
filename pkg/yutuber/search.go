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

func WithSearchChannelId(channelId string) searchOption {
	return func(s *search) {
		if channelId != "" {
			s.channelId = channelId
		}
	}
}

func WithSearchChannelType(channelType string) searchOption {
	return func(s *search) {
		if channelType != "" {
			s.channelType = channelType
		}
	}
}

func WithSearchEventType(eventType string) searchOption {
	return func(s *search) {
		if eventType != "" {
			s.eventType = eventType
		}
	}
}

func WithSearchForContentOwner(forContentOwner bool) searchOption {
	return func(s *search) {
		s.forContentOwner = forContentOwner
	}
}

func WithSearchForDeveloper(forDeveloper bool) searchOption {
	return func(s *search) {
		s.forDeveloper = forDeveloper
	}
}

func WithSearchForMine(forMine bool) searchOption {
	return func(s *search) {
		s.forMine = forMine
	}
}

func WithSearchLocation(location string) searchOption {
	return func(s *search) {
		if location != "" {
			s.location = location
		}
	}
}

func WithSearchLocationRadius(locationRadius string) searchOption {
	return func(s *search) {
		if locationRadius != "" {
			s.locationRadius = locationRadius
		}
	}
}

func WithSearchMaxResults(maxResults int64) searchOption {
	return func(s *search) {
		s.maxResults = maxResults
	}
}

func WithSearchOnBehalfOfContentOwner(onBehalfOfContentOwner string) searchOption {
	return func(s *search) {
		if onBehalfOfContentOwner != "" {
			s.onBehalfOfContentOwner = onBehalfOfContentOwner
		}
	}
}

func WithSearchOrder(order string) searchOption {
	return func(s *search) {
		if order != "" {
			s.order = order
		}
	}
}

func WithSearchPublishedAfter(publishedAfter string) searchOption {
	return func(s *search) {
		if publishedAfter != "" {
			s.publishedAfter = publishedAfter
		}
	}
}

func WithSearchPublishedBefore(publishedBefore string) searchOption {
	return func(s *search) {
		if publishedBefore != "" {
			s.publishedBefore = publishedBefore
		}
	}
}

func WithSearchQ(q string) searchOption {
	return func(s *search) {
		if q != "" {
			s.q = q
		}
	}
}

func WithSearchRegionCode(regionCode string) searchOption {
	return func(s *search) {
		if regionCode != "" {
			s.regionCode = regionCode
		}
	}
}

func WithSearchRelevanceLanguage(relevanceLanguage string) searchOption {
	return func(s *search) {
		if relevanceLanguage != "" {
			s.relevanceLanguage = relevanceLanguage
		}
	}
}

func WithSearchSafeSearch(safeSearch string) searchOption {
	return func(s *search) {
		if safeSearch != "" {
			s.safeSearch = safeSearch
		}
	}
}

func WithSearchTopicId(topicId string) searchOption {
	return func(s *search) {
		if topicId != "" {
			s.topicId = topicId
		}
	}
}

func WithSearchTypes(types string) searchOption {
	return func(s *search) {
		if types != "" {
			s.types = types
		}
	}
}

func WithSearchVideoCaption(videoCaption string) searchOption {
	return func(s *search) {
		if videoCaption != "" {
			s.videoCaption = videoCaption
		}
	}
}

func WithSearchVideoCategoryId(videoCategoryId string) searchOption {
	return func(s *search) {
		if videoCategoryId != "" {
			s.videoCategoryId = videoCategoryId
		}
	}
}

func WithSearchVideoDefinition(videoDefinition string) searchOption {
	return func(s *search) {
		if videoDefinition != "" {
			s.videoDefinition = videoDefinition
		}
	}
}

func WithSearchVideoDimension(videoDimension string) searchOption {
	return func(s *search) {
		if videoDimension != "" {
			s.videoDimension = videoDimension
		}
	}
}

func WithSearchVideoDuration(videoDuration string) searchOption {
	return func(s *search) {
		if videoDuration != "" {
			s.videoDuration = videoDuration
		}
	}
}

func WithSearchVideoEmbeddable(videoEmbeddable string) searchOption {
	return func(s *search) {
		if videoEmbeddable != "" {
			s.videoEmbeddable = videoEmbeddable
		}
	}
}

func WithSearchVideoLicense(videoLicense string) searchOption {
	return func(s *search) {
		if videoLicense != "" {
			s.videoLicense = videoLicense
		}
	}
}

func WithSearchVideoPaidProductPlacement(videoPaidProductPlacement string) searchOption {
	return func(s *search) {
		if videoPaidProductPlacement != "" {
			s.videoPaidProductPlacement = videoPaidProductPlacement
		}
	}
}

func WithSearchVideoSyndicated(videoSyndicated string) searchOption {
	return func(s *search) {
		if videoSyndicated != "" {
			s.videoSyndicated = videoSyndicated
		}
	}
}

func WithSearchVideoType(videoType string) searchOption {
	return func(s *search) {
		if videoType != "" {
			s.videoType = videoType
		}
	}
}
