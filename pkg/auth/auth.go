package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	cacheTokenFile = "youtube.token.json"
)

var (
	state            = utils.RandomStage()
	errStateMismatch = errors.New("state doesn't match, possible CSRF attack")
	errReadPrompt    = errors.New("unable to read prompt")
	errExchange      = errors.New("unable retrieve token from web or prompt")
	errStartWeb      = errors.New("unable to start web server")
	errCacheToken    = errors.New("unable to cache token")
	errParseToken    = errors.New("unable to parse token")
	errRefreshToken  = errors.New("unable to refresh token, please re-authenticate in cli")
	errParseSecret   = errors.New("unable to parse client secret")

	scope = []string{
		youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
		youtube.YoutubeChannelMembershipsCreatorScope,
	}
)

const manualInputHint = `
After completing the authorization flow, enter the authorization code on command line.

If you end up in an error page after completing the authorization flow,
and the url in the address bar is in the form of
'localhost:8216/?state=DONOT-COPY&code=COPY-THIS&scope=DONOT-COPY'
ONLY 'COPY-THIS' part is the code you need to enter on command line.
`

func InitClient(
	ctx context.Context, cred, token string, cacheable bool,
) *http.Client {
	config := getConfig(cred, scope...)

	authedToken := &oauth2.Token{}
	err := json.Unmarshal([]byte(token), authedToken)
	if err != nil {
		log.Fatalln(errors.Join(errParseToken, err))
	}

	if !authedToken.Valid() {
		tokenSource := config.TokenSource(ctx, authedToken)
		authedToken, err = tokenSource.Token()
		if err != nil && cacheable {
			var client *http.Client
			client, authedToken = newClient(ctx, config)
			saveToken(cacheTokenFile, authedToken)
			return client
		} else if err != nil {
			log.Fatalln(errors.Join(errRefreshToken, err))
		}

		if authedToken != nil && cacheable {
			saveToken(cacheTokenFile, authedToken)
		}
		return config.Client(ctx, authedToken)
	}

	return config.Client(ctx, authedToken)
}

func newClient(
	ctx context.Context, config *oauth2.Config,
) (client *http.Client, token *oauth2.Token) {
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	token = getTokenFromWeb(config, authURL)
	client = config.Client(ctx, token)

	return
}

func getConfig(cred string, scope ...string) *oauth2.Config {
	config, err := google.ConfigFromJSON([]byte(cred), scope...)
	if err != nil {
		log.Fatalln(errors.Join(errParseSecret, err))
	}

	return config
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
				_ = listener.Close()
				w.Header().Set("Content-Type", "text/plain")
				_, _ = fmt.Fprintf(
					w, "Received code: %s\r\nYou can now safely close this window.", code,
				)
			},
		),
	)

	return codeCh
}

func getCodeFromPrompt(authURL string) (code string) {
	fmt.Printf(
		"It seems that your browser is not open. Go to the following "+
			"link in your browser:\n%s\n", authURL,
	)
	fmt.Print(manualInputHint)
	_, err := fmt.Scan(&code)
	if err != nil {
		log.Fatalln(errors.Join(errReadPrompt, err))
	}

	if strings.HasPrefix(code, "4%2F") {
		code = strings.Replace(code, "4%2F", "4/", 1)
	}
	return code
}

func getTokenFromWeb(config *oauth2.Config, authURL string) *oauth2.Token {
	codeCh := startWebServer(config.RedirectURL)

	var code string
	if err := utils.OpenURL(authURL); err == nil {
		fmt.Printf(
			"Your browser has been opened to an authorization URL. This "+
				"program will resume once authorization has been provided.\n%s\n",
			authURL,
		)
		code = <-codeCh
	}

	if code == "" {
		code = getCodeFromPrompt(authURL)
	}

	fmt.Printf("Authorization code generated: %s\n", code)
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalln(errors.Join(errExchange, err))
	}

	return token
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln(errors.Join(errCacheToken, err))
	}

	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalln(errors.Join(errCacheToken, err))
	}
}
