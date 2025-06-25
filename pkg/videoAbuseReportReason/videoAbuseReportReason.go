package videoAbuseReportReason

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"io"
)

var (
	service                      *youtube.Service
	errGetVideoAbuseReportReason = errors.New("failed to get video abuse report reason")
)

type videoAbuseReportReason struct {
	Hl string `yaml:"hl" json:"hl"`
}

type VideoAbuseReportReason interface {
	Get([]string) ([]*youtube.VideoAbuseReportReason, error)
	List([]string, string, io.Writer) error
}

type Option func(*videoAbuseReportReason)

func NewVideoAbuseReportReason(opt ...Option) VideoAbuseReportReason {
	va := &videoAbuseReportReason{}
	for _, o := range opt {
		o(va)
	}
	return va
}

func (va *videoAbuseReportReason) Get(parts []string) (
	[]*youtube.VideoAbuseReportReason, error,
) {
	call := service.VideoAbuseReportReasons.List(parts)
	if va.Hl != "" {
		call = call.Hl(va.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoAbuseReportReason, err)
	}

	return res.Items, nil
}

func (va *videoAbuseReportReason) List(
	parts []string, output string, writer io.Writer,
) error {
	videoAbuseReportReasons, err := va.Get(parts)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		utils.PrintJSON(videoAbuseReportReasons, writer)
	case "yaml":
		utils.PrintYAML(videoAbuseReportReasons, writer)
	default:
		_, _ = fmt.Fprintln(writer, "ID\tTitle")
		for _, videoAbuseReportReason := range videoAbuseReportReasons {
			_, _ = fmt.Fprintf(
				writer, "%s\t%s\n",
				videoAbuseReportReason.Id, videoAbuseReportReason.Snippet.Label,
			)
		}
	}
	return nil
}

func WithHL(hl string) Option {
	return func(vc *videoAbuseReportReason) {
		vc.Hl = hl
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *videoAbuseReportReason) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}
