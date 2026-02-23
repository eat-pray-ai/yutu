package common

import (
	"fmt"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
)

type Fields struct {
	Service  *youtube.Service `yaml:"-" json:"-"`
	Parts    []string         `yaml:"parts" json:"parts,omitempty"`
	Output   string           `yaml:"output" json:"output,omitempty"`
	Jsonpath string           `yaml:"jsonpath" json:"jsonpath,omitempty"`
}

func (d *Fields) GetFields() *Fields {
	return d
}

func (d *Fields) EnsureService() {
	if d.Service == nil {
		svc, err := auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
		if err != nil {
			panic(fmt.Sprintf("failed to create YouTube service: %v", err))
		}
		d.Service = svc
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
