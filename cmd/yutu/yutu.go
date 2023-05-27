package main

import (
	"context"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/you2be"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	client := auth.GetClient(youtube.YoutubeReadonlyScope)
	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	auth.HandleError(err, "Error creating YouTube client")

	you2be.ChannelsListById(service, []string{"snippet", "contentDetails", "statistics"}, "***REMOVED***")
}
