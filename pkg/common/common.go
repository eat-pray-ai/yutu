package common

import (
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
)

type Fields struct {
	Service  *youtube.Service
	Parts    []string `yaml:"parts" json:"parts"`
	Output   string   `yaml:"output" json:"output"`
	Jsonpath string   `yaml:"jsonpath" json:"jsonpath"`
}

func (d *Fields) GetFields() *Fields {
	return d
}

func (d *Fields) EnsureService() {
	if d.Service == nil {
		d.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

type HasFields interface {
	GetFields() *Fields
	EnsureService()
}

func WithParts[T HasFields](parts []string) func(T) {
	return func(t T) {
		t.GetFields().Parts = parts
	}
}

func WithOutput[T HasFields](output string) func(T) {
	return func(t T) {
		t.GetFields().Output = output
	}
}

func WithJsonpath[T HasFields](jsonpath string) func(T) {
	return func(t T) {
		t.GetFields().Jsonpath = jsonpath
	}
}

func WithService[T HasFields](svc *youtube.Service) func(T) {
	return func(t T) {
		t.GetFields().Service = svc
	}
}
