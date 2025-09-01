package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

var (
	errCreateSvc  = errors.New("unable to create YouTube service")
	errReadToken  = errors.New("unable to read token")
	errReadSecret = errors.New("unable to read client secret")
)

const missingClientSecretsHint string = `
Please configure OAuth 2.0
You need to populate the client_secrets.json file
found at: %s
with information from the {{ Google Cloud Console }}
{{ https://cloud.google.com/console }}
For more information about the client_secrets.json file format, please visit:
https://developers.google.com/api-client-library/python/guide/aaa_client_secrets
`

type svc struct {
	Cacheable  bool   `yaml:"cacheable" json:"cacheable"`
	Credential string `yaml:"credential" json:"credential"`
	CacheToken string `yaml:"cache_token" json:"cache_token"`
	service    *youtube.Service
	ctx        context.Context
}

type Svc interface {
	GetService() *youtube.Service
	refreshClient() *http.Client
	newClient(*oauth2.Config) (*http.Client, *oauth2.Token)
	getConfig() *oauth2.Config
	startWebServer(string) chan string
	getTokenFromWeb(*oauth2.Config, string) *oauth2.Token
	getCodeFromPrompt(string) string
	saveToken(string, *oauth2.Token)
}

type Option func(*svc)

func NewY2BService(opts ...Option) Svc {
	s := &svc{}
	s.ctx = context.Background()

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithCredential(cred string) Option {
	return func(s *svc) {
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
				fmt.Printf(missingClientSecretsHint, cred)
				log.Fatalln(errors.Join(errReadSecret, err))
			}
			s.Credential = string(credBytes)
			return
		}

		if credB64, err := base64.StdEncoding.DecodeString(cred); err == nil {
			s.Credential = string(credB64)
		} else if utils.IsJson(cred) {
			s.Credential = cred
		} else {
			fmt.Printf(missingClientSecretsHint, cred)
			log.Fatalln(errors.Join(errReadSecret, err))
		}
	}
}

func WithCacheToken(token string) Option {
	return func(s *svc) {
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
			s.CacheToken = string(tokenBytes)
			s.Cacheable = true
			return
		} else if os.IsNotExist(err) && strings.HasSuffix(token, ".json") {
			s.CacheToken = token
			s.Cacheable = true
			return
		}

		if tokenB64, err := base64.StdEncoding.DecodeString(token); err == nil {
			s.CacheToken = string(tokenB64)
		} else if utils.IsJson(token) {
			s.CacheToken = token
		} else {
			log.Fatalln(errors.Join(errReadToken, err))
		}
	}
}
