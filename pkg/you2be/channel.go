package you2be

import (
	"fmt"

	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/auth"
)

func ChannelsListById(service *youtube.Service, part []string, channelId string) {
	call := service.Channels.List(part)
	call = call.Id(channelId)
	response, err := call.Do()
	auth.HandleError(err, "")

	fmt.Printf("This channel's ID is %s.\nIts title is '%s', and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount)
}
