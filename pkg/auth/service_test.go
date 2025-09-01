package auth

import (
	"context"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestNewY2BService(t *testing.T) {
	credential := `{"client_id":"test"}`
	cacheToken := `{"access_token":"test"}`

	mockFS := fstest.MapFS{
		"youtube.token.json": &fstest.MapFile{
			Data: []byte(cacheToken),
		},
		"client_secrets.json": &fstest.MapFile{
			Data: []byte(credential),
		},
	}

	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want Svc
	}{
		{
			name: "with all options - json",
			args: args{
				opts: []Option{
					WithCredential(credential, mockFS),
					WithCacheToken(cacheToken, mockFS),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				Cacheable:  false,
				ctx:        context.Background(),
			},
		},
		{
			name: "with all options - base64",
			args: args{
				opts: []Option{
					WithCredential(
						"eyJjbGllbnRfaWQiOiJ0ZXN0In0=", mockFS,
					),
					WithCacheToken(
						"eyJhY2Nlc3NfdG9rZW4iOiJ0ZXN0In0=", mockFS,
					),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				Cacheable:  false,
				ctx:        context.Background(),
			},
		},
		{
			name: "with all options - file",
			args: args{
				opts: []Option{
					WithCredential("/client_secrets.json", mockFS),
					WithCacheToken("/youtube.token.json", mockFS),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				Cacheable:  true,
				ctx:        context.Background(),
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &svc{
				ctx: context.Background(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewY2BService(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewY2BService() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
