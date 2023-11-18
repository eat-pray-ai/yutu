package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/eat-pray-ai/yutu/pkg/util"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	authCmd    = flag.NewFlagSet("auth", flag.ExitOnError)
	credential = authCmd.String("file", "client_secret.json", "credential filename")
	scope      = authCmd.String("scope", "youtube.readonly", "scope to authenticate")

	channelCmd = flag.NewFlagSet("channel", flag.ExitOnError)
	cId        = channelCmd.String("id", "", "channel ID")
	cList      = channelCmd.Bool("list", false, "list channel")

	videoCmd = flag.NewFlagSet("video", flag.ExitOnError)
	filename = videoCmd.String("file", "", "video file to upload")
	title    = videoCmd.String("title", "", "video title")
	desc     = videoCmd.String("desc", "", "video description")
	category = videoCmd.String("category", "", "video category")
	keywords = videoCmd.String("keywords", "", "comma separated of video keywords")
	privacy  = videoCmd.String("privacy", "unlisted", "video privacy status")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <command> [options]", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), `
Commands:
  auth:
    -file <credential filename>
    -scope <scope to authenticate>
  channel:
    -id <channel id>
    -list
  video:
    -file <video file to upload>
    -title <video title>
    -desc <video description>
    -category <video category>
    -keywords <comma separated of video keywords>
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

		if *cList {
			channel := yutuber.ChannelsGet(service, *cId)
			yutuber.ChannelsList(channel)
		}
	case "video":
		videoCmd.Parse(os.Args[2:])
		client := auth.GetClient(ctx, "", youtube.YoutubeUploadScope)
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		util.HandleError(err, "Error creating YouTube client")

		snippet := &youtube.VideoSnippet{
			Title:       *title,
			Description: *desc,
			CategoryId:  *category,
		}

		yutuber.VideoInsert(service, *filename, snippet, *keywords, *privacy)
	}
}
