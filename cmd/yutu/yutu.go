package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/you2be"

	"github.com/eat-pray-ai/yutu/pkg/util"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	authCmd    = flag.NewFlagSet("auth", flag.ExitOnError)
	credential = authCmd.String("cred", "client_secret.json", "Credential file")
	scope      = authCmd.String("scope", "youtube.readonly", "Scope to authenticate")

	channelCmd   = flag.NewFlagSet("channel", flag.ExitOnError)
	channelId    = channelCmd.String("id", "", "Channel ID")
	channelsList = channelCmd.Bool("list", false, "List channel")

	videoCmd    = flag.NewFlagSet("video", flag.ExitOnError)
	filename    = videoCmd.String("file", "", "Name of video file to upload")
	title       = videoCmd.String("title", "", "Video title")
	description = videoCmd.String("description", "", "Video description")
	category    = videoCmd.String("category", "", "Video category")
	keywords    = videoCmd.String("keywords", "", "Comma separated list of video keywords")
	privacy     = videoCmd.String("privacy", "unlisted", "Video privacy status")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <command> [options]", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), `
Commands:
  auth:
    -cred <credential file>
    -scope <scope to authenticate>
  channel:
    -id <channel ID>
    -list
  video:
    -file <video file>
    -title <video title>
    -description <video description>
    -category <video category>
    -keywords <video keywords>
    -privacy <video privacy status>
`)

		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	fmt.Printf("Args: %v\n", os.Args[1:])

	switch os.Args[1] {
	case "auth":
		authCmd.Parse(os.Args[2:])
		auth.GetClient(ctx, *credential, *scope)
	case "channel":
		channelCmd.Parse(os.Args[2:])
		client := auth.GetClient(ctx, "", youtube.YoutubeReadonlyScope)
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		util.HandleError(err, "Error creating YouTube client")

		if *channelsList {
			channel := you2be.ChannelsGet(service, *channelId)
			you2be.ChannelsList(channel)
		}
	case "video":
		videoCmd.Parse(os.Args[2:])
		client := auth.GetClient(ctx, "", youtube.YoutubeUploadScope)
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		util.HandleError(err, "Error creating YouTube client")

		you2be.VideoInsert(service, *filename, *title, *description, *category, *keywords, *privacy)
	}
}
