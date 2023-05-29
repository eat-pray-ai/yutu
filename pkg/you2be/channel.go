package you2be

import (
	"fmt"

	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/auth"
)

var part = []string{"snippet", "statistics"}

func ChannelsListById(service *youtube.Service, channelId string) *youtube.Channel {
	call := service.Channels.List(part)
	call = call.Id(channelId)
	response, err := call.Do()
	auth.HandleError(err, "")

	channel := response.Items[0]
	fmt.Printf("Channel ID: %s\nTitle: %s\nDescription: %s\nPublished At: %s\nCountry: %s\nSubscriber count: %d\nVideo count: %d\nViews: %d\n",
		channel.Id,
		channel.Snippet.Title,
		channel.Snippet.Description,
		channel.Snippet.PublishedAt,
		channel.Snippet.Country,
		channel.Statistics.SubscriberCount,
		channel.Statistics.VideoCount,
		channel.Statistics.ViewCount)

	return response.Items[0]
}

func ChannelsUpdate(service *youtube.Service, channelId string) {
	channel := ChannelsListById(service, channelId)
	channel.Snippet.Title = "看剧啦饭酱"
	call := service.Channels.Update(part, channel)
	channel, err := call.Do()
	auth.HandleError(err, "")

	fmt.Printf("Channel ID: %s\nTitle: %s\nDescription: %s\nPublished At: %s\nCountry: %s\nSubscriber count: %d\nVideo count: %d\nViews: %d\n",
		channel.Id,
		channel.Snippet.Title,
		channel.Snippet.Description,
		channel.Snippet.PublishedAt,
		channel.Snippet.Country,
		channel.Statistics.SubscriberCount,
		channel.Statistics.VideoCount,
		channel.Statistics.ViewCount)
}
