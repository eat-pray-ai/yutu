package you2be

import (
	"fmt"

	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/util"
)

var part = []string{"snippet", "statistics"}

func ChannelsGet(service *youtube.Service, channelId string) *youtube.Channel {
	call := service.Channels.List(part)
	call = call.Id(channelId)
	response, err := call.Do()
	util.HandleError(err, "")

	return response.Items[0]
}

func ChannelsList(channel *youtube.Channel) {
	fmt.Printf(`Channel ID: %s
Title: %s
Description: %s
Published At: %s
Country: %s
Subscriber count: %d
Video count: %d
Views: %d`,
		channel.Id,
		channel.Snippet.Title,
		channel.Snippet.Description,
		channel.Snippet.PublishedAt,
		channel.Snippet.Country,
		channel.Statistics.SubscriberCount,
		channel.Statistics.VideoCount,
		channel.Statistics.ViewCount)
}

func ChannelsUpdate(service *youtube.Service, channelId string, title string) {
	channel := ChannelsGet(service, channelId)
	channel.Snippet.Title = title
	call := service.Channels.Update(part, channel)
	channel, err := call.Do()
	util.HandleError(err, "")

	ChannelsList(channel)
}
