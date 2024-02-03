package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/util"
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
)

var (
	authCmd    = flag.NewFlagSet("auth", flag.ExitOnError)
	credential = authCmd.String("file", "client_secret.json", "credential filename")
	scope      = authCmd.String("scope", youtube.YoutubeReadonlyScope, "scope to authenticate")

	channelCmd = flag.NewFlagSet("channel", flag.ExitOnError)
	cId        = channelCmd.String("id", "", "channel ID")
	cList      = channelCmd.Bool("list", false, "list channel")

	videoCmd = flag.NewFlagSet("video", flag.ExitOnError)
	path     = videoCmd.String("path", "", "path to video file for uploading")
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
    -path <path to video file for uploading>
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

	switch os.Args[1] {
	case "auth":
		authCmd.Parse(os.Args[2:])
		auth.GetClient(ctx, *credential, *scope)
	case "channel":
		channelCmd.Parse(os.Args[2:])
		client := auth.GetClient(ctx, "client_secret.json", youtube.YoutubeReadonlyScope)
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		util.HandleError(err, "Error creating YouTube client")

		if *cList {
			channel := yutuber.ChannelsGet(service, *cId)
			yutuber.ChannelsList(channel)
		}
	case "video":
		videoCmd.Parse(os.Args[2:])
		client := auth.GetClient(ctx, "client_secret.json", youtube.YoutubeUploadScope)
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		util.HandleError(err, "Error creating YouTube client")

		newVideo := &yutuber.Video{
			Path:     *path,
			Title:    *title,
			Desc:     *desc,
			Category: *category,
			Keywords: *keywords,
			Privacy:  *privacy,
		}

		newVideo.Insert(service)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
