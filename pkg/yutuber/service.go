package yutuber

import (
	"context"
	"log"

	"github.com/eat-pray-ai/yutu/pkg/auth"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type authService struct {
	Scope string
}

type AuthService interface {
	auth(context.Context) *youtube.Service
}

func (as *authService) auth(ctx context.Context) *youtube.Service {
	client := auth.GetClient(ctx, "client_secret.json", as.Scope)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	return service
}
