package main

import (
	"context"
	"log"
	"os"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/you2be"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	ctx := context.Background()

	b, err := os.ReadFile("client_secret_desktop.json")
	if err != nil {
		log.Fatalf("Unable to read client secret: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file %s to config: %v", string(b), err)
	}

	client := auth.GetClient(ctx, config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	auth.HandleError(err, "Error creating YouTube client")

	you2be.ChannelsListById(service, []string{"snippet", "contentDetails", "statistics"}, "***REMOVED***")
}
