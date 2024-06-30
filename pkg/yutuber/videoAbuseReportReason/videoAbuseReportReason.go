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
	Hl string `yaml:"hl" json:"hl"`
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

func (va *videoAbuseReportReason) get(parts []string) []*youtube.VideoAbuseReportReason {
	call := service.VideoAbuseReportReasons.List(parts)
	if va.Hl != "" {
		call = call.Hl(va.Hl)
	}

	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(va)
		log.Fatalln(errors.Join(errGetVideoAbuseReportReason, err))
	}

	return res.Items
}

func (va *videoAbuseReportReason) List(parts []string, output string) {
	videoAbuseReportReasons := va.get(parts)
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
		vc.Hl = hl
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
