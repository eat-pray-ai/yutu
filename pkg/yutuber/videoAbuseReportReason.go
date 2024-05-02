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
	errGetVideoAbuseReportReason error = errors.New("failed to get video abuse report reason")
)

type videoAbuseReportReason struct{}

type VideoAbuseReportReason interface {
	get([]string) []*youtube.VideoAbuseReportReason
	List([]string, string)
}

type videoAbuseReportReasonOption func(*videoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...videoAbuseReportReasonOption) VideoAbuseReportReason {
	service = auth.NewY2BService()
	va := &videoAbuseReportReason{}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (vc *videoAbuseReportReason) get(parts []string) []*youtube.VideoAbuseReportReason {
	call := service.VideoAbuseReportReasons.List(parts)
	response, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideoAbuseReportReason, err))
	}

	return response.Items
}

func (vc *videoAbuseReportReason) List(parts []string, output string) {
	videoAbuseReportReasons := vc.get(parts)
	switch output {
	case "json":
		utils.PrintJSON(videoAbuseReportReasons)
	case "yaml":
		utils.PrintYAML(videoAbuseReportReasons)
	default:
		fmt.Println("ID\tTitle")
		for _, videoAbuseReportReason := range videoAbuseReportReasons {
			fmt.Printf(
				"%s\t%s\n", videoAbuseReportReason.Id,
				videoAbuseReportReason.Snippet.Label,
			)
		}
	}
}
