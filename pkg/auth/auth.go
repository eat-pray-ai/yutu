package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

var (
	state            = utils.RandomStage()
	credential       = "client_secret.json"
	cacheToken       = "youtube.token.json"
	errStateMismatch = errors.New("state doesn't match, possible CSRF attack")
	errCreateSvc     = errors.New("unable to create YouTube service")
	errReadPrompt    = errors.New("unable to read prompt")
	errExchange      = errors.New("unable retrieve token from web or prompt")
	errStartWeb      = errors.New("unable to start web server")
	errCacheToken    = errors.New("unable to cache token")
	errReadSecret    = errors.New("unable to read client secret file")
	errParseSecret   = errors.New("unable to parse client secret to config")
)

const missingClientSecretsMessage string = `
Please configure OAuth 2.0
To make this sample run, you need to populate the client_secrets.json file
found at:
  %v
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
	scope := []string{
		youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
		youtube.YoutubeChannelMembershipsCreatorScope,
	}
	client := getClient(ctx, scope...)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(errors.Join(errCreateSvc, err))
	}

	return service
}

func getClient(ctx context.Context, scope ...string) *http.Client {
	config := getConfig(scope...)

	token, err := tokenFromCache(cacheToken)
	if err != nil {
		return newClient(ctx, config, cacheToken)
	} else if !token.Valid() {
		tokenSource := config.TokenSource(ctx, token)
		token, err = tokenSource.Token()
		if token != nil {
			saveToken(cacheToken, token)
		}
		if err != nil {
			return newClient(ctx, config, cacheToken)
		}
	}

	return config.Client(ctx, token)
}

func getConfig(scope ...string) *oauth2.Config {
	cred, err := os.ReadFile(credential)
	if err != nil {
		fmt.Printf(missingClientSecretsMessage, credential)
		log.Fatalln(errors.Join(errReadSecret, err))
	}

	config, err := google.ConfigFromJSON(cred, scope...)
	if err != nil {
		log.Fatalln(errors.Join(errParseSecret, err))
	}

	return config
}

func newClient(ctx context.Context, config *oauth2.Config, cacheToken string) *http.Client {
	var token *oauth2.Token
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	token = getTokenFromWeb(config, authURL)
	saveToken(cacheToken, token)

	return config.Client(ctx, token)
}

func startWebServer(redirectURL string) chan string {
	u, err := url.Parse(redirectURL)
	if err != nil {
		log.Fatalln(errors.Join(errStartWeb, err))
	}

	listener, err := net.Listen("tcp", u.Host)
	if err != nil {
		log.Fatalln(errors.Join(errStartWeb, err))
	}

	codeCh := make(chan string)
	go http.Serve(
		listener, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					return
				}
				s := r.FormValue("state")
				if s != state {
					log.Fatalf("%v: %s != %s\n", errStateMismatch, s, state)
				}
				code := r.FormValue("code")
				codeCh <- code
				listener.Close()
				w.Header().Set("Content-Type", "text/plain")
				fmt.Fprintf(
					w, "Received code: %v\r\nYou can now safely close this window.", code,
				)
			},
		),
	)

	return codeCh
}

func getCodeFromPrompt(authURL string) string {
	var code string
	fmt.Printf(
		"It seems that your browser is not open. Go to the following "+
			"link in your browser:\n%v\n", authURL,
	)
	fmt.Print(
		"After completing the authorization flow, enter the authorization " +
			"code on command line. \nIf you end up in an error page after completing " +
			"the authorization flow, and the url in the address bar is in the form of " +
			"\n'localhost:8216/?state=DONOT-COPY&code=COPY-THIS&scope=DONOT-COPY'\n" +
			"ONLY 'COPY-THIS' is the code you need to enter on command line.\n",
	)
	_, err := fmt.Scan(&code)
	if err != nil {
		log.Fatalln(errors.Join(errReadPrompt, err))
	}

	return code
}

func getTokenFromWeb(config *oauth2.Config, authURL string) *oauth2.Token {
	codeCh := startWebServer(config.RedirectURL)

	var code string
	if err := utils.OpenURL(authURL); err == nil {
		fmt.Printf(
			"Your browser has been opened to an authorization URL. This "+
				"program will resume once authorization has been provided.\n%v\n",
			authURL,
		)
		code = <-codeCh
	}

	if code == "" {
		code = getCodeFromPrompt(authURL)
	}

	fmt.Printf("Authorization code generated: %v\n", code)
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalln(errors.Join(errExchange, err))
	}

	return token
}

func tokenFromCache(cacheToken string) (*oauth2.Token, error) {
	f, err := os.Open(cacheToken)
	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(cacheToken string, token *oauth2.Token) {
	f, err := os.OpenFile(cacheToken, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln(errors.Join(errCacheToken, err))
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func WithCredential(cred string) Option {
	return func() {
		credential = cred
	}
}

func WithCacheToken(token string) Option {
	return func() {
		cacheToken = token
	}
}
