package main

import (
	"context"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/you2be"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	ctx := context.Background()
	client := auth.GetClient(ctx, youtube.YoutubeReadonlyScope)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	auth.HandleError(err, "Error creating YouTube client")

	you2be.ChannelsListById(service, "***REMOVED***")
}
