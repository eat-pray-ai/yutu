package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
	"strings"
)

var (
	cacheable     bool
	credential    string
	cacheToken    string
	errCreateSvc  = errors.New("unable to create YouTube service")
	errReadToken  = errors.New("unable to read token")
	errReadSecret = errors.New("unable to read client secret")
)

const missingClientSecretsHint string = `
Please configure OAuth 2.0
To make this sample run, you need to populate the client_secrets.json file
found at:
  %s
with information from the {{ Google Cloud Console }}
{{ https://cloud.google.com/console }}
For more information about the client_secrets.json file format, please visit:
https://developers.google.com/api-client-library/python/guide/aaa_client_secrets
`

type Option func()

func NewY2BService(opts ...Option) *youtube.Service {
	for _, opt := range opts {
		opt()
	}

	ctx := context.Background()
	client := InitClient(ctx, credential, cacheToken, cacheable)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(errors.Join(errCreateSvc, err))
	}

	return service
}

func WithCredential(cred string) Option {
	return func() {
		// cred > YUTU_CREDENTIAL
		envCred, ok := os.LookupEnv("YUTU_CREDENTIAL")
		if cred == "" && ok {
			cred = envCred
		} else if cred == "" {
			cred = "client_secret.json"
		}

		// 1. cred is a file path
		// 2. cred is a base64 encoded string
		// 3. cred is a json string
		if _, err := os.Stat(cred); err == nil {
			credBytes, err := os.ReadFile(cred)
			if err != nil {
				fmt.Printf(missingClientSecretsHint, credential)
				log.Fatalln(errors.Join(errReadSecret, err))
			}
			credential = string(credBytes)
			return
		}

		if credB64, err := base64.StdEncoding.DecodeString(cred); err == nil {
			credential = string(credB64)
		} else if utils.IsJson(cred) {
			credential = cred
		} else {
			fmt.Printf(missingClientSecretsHint, credential)
			log.Fatalln(errors.Join(errReadSecret, err))
		}
	}
}

func WithCacheToken(token string) Option {
	return func() {
		// token > YUTU_CACHE_TOKEN
		envToken, ok := os.LookupEnv("YUTU_CACHE_TOKEN")
		if token == "" && ok {
			token = envToken
		} else if token == "" {
			token = cacheTokenFile
		}

		// 1. token is a file path
		// 2. token is a base64 encoded string
		// 3. token is a json string
		if _, err := os.Stat(token); err == nil {
			tokenBytes, err := os.ReadFile(token)
			if err != nil {
				log.Fatalln(errors.Join(errReadToken, err))
			}
			cacheToken = string(tokenBytes)
			cacheable = true
			return
		} else if os.IsNotExist(err) && strings.HasSuffix(token, ".json") {
			cacheToken = token
			cacheable = true
			return
		}

		if tokenB64, err := base64.StdEncoding.DecodeString(token); err == nil {
			cacheToken = string(tokenB64)
		} else if utils.IsJson(token) {
			cacheToken = token
		} else {
			log.Fatalln(errors.Join(errReadToken, err))
		}
	}
}
