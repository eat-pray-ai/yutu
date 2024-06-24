package videoAbuseReportReason

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
)

var (
	service                      *youtube.Service
	errGetVideoAbuseReportReason = errors.New("failed to get video abuse report reason")
)

type videoAbuseReportReason struct {
	hl string
}

type VideoAbuseReportReason interface {
	get([]string) []*youtube.VideoAbuseReportReason
	List([]string, string)
}

type Option func(*videoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) VideoAbuseReportReason {
	va := &videoAbuseReportReason{}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (vc *videoAbuseReportReason) get(parts []string) []*youtube.VideoAbuseReportReason {
	call := service.VideoAbuseReportReasons.List(parts)
	if vc.hl != "" {
		call = call.Hl(vc.hl)
	}

	res, err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errGetVideoAbuseReportReason, err))
	}

	return res.Items
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

func WithHL(hl string) Option {
	return func(vc *videoAbuseReportReason) {
		vc.hl = hl
	}
}

func WithService(svc *youtube.Service) Option {
	return func(vc *videoAbuseReportReason) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}
