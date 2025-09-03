package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	createSvcFailed    = "failed to create YouTube service"
	parseUrlFailed     = "failed to parse redirect URL"
	stateMatchFailed   = "state doesn't match, possible CSRF attack"
	readPromptFailed   = "failed to read prompt"
	exchangeFailed     = "failed to exchange token"
	listenFailed       = "failed to start web server"
	cacheTokenFailed   = "failed to cache token"
	parseTokenFailed   = "failed to parse token"
	refreshTokenFailed = "failed to refresh token, please re-authenticate in cli"
	parseSecretFailed  = "failed to parse client secret"

	browserOpenedHint = "Your browser has been opened to an authorization URL. yutu will resume once authorization has been provided.\n%s\n"
	openBrowserHint   = "It seems that your browser is not open. Go to the following link in your browser:\n%s\n"
	manualInputHint   = `
After completing the authorization flow, enter the authorization code on command line.

If you end up in an error page after completing the authorization flow,
and the url in the address bar is in the form of
'localhost:8216/?state=DONOT-COPY&code=COPY-THIS&scope=DONOT-COPY'
ONLY 'COPY-THIS' part is the code you need to enter on command line.
`
)

var (
	state = utils.RandomStage()
	scope = []string{
		youtube.YoutubeScope,
		youtube.YoutubeForceSslScope,
		youtube.YoutubeChannelMembershipsCreatorScope,
	}
)

func (s *svc) GetService() *youtube.Service {
	client := s.refreshClient()
	service, err := youtube.NewService(s.ctx, option.WithHTTPClient(client))
	if err != nil {
		slog.Error(createSvcFailed, "error", err)
		os.Exit(1)
	}
	s.service = service

	return s.service
}

func (s *svc) refreshClient() (client *http.Client) {
	config := s.getConfig()
	authedToken := &oauth2.Token{}
	err := json.Unmarshal([]byte(s.CacheToken), authedToken)
	if err != nil {
		client, authedToken = s.newClient(config)
		if s.tokenFile != "" {
			s.saveToken(authedToken)
		}
		return client
	}

	if !authedToken.Valid() {
		tokenSource := config.TokenSource(s.ctx, authedToken)
		authedToken, err = tokenSource.Token()
		if err != nil && s.tokenFile != "" {
			client, authedToken = s.newClient(config)
			s.saveToken(authedToken)
			return client
		} else if err != nil {
			slog.Error(refreshTokenFailed, "error", err)
			os.Exit(1)
		}

		if authedToken != nil && s.tokenFile != "" {
			s.saveToken(authedToken)
		}
		return config.Client(s.ctx, authedToken)
	}

	return config.Client(s.ctx, authedToken)
}

func (s *svc) newClient(config *oauth2.Config) (
	client *http.Client, token *oauth2.Token,
) {
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	token = s.getTokenFromWeb(config, authURL)
	client = config.Client(s.ctx, token)

	return
}

func (s *svc) getConfig() *oauth2.Config {
	config, err := google.ConfigFromJSON([]byte(s.Credential), scope...)
	if err != nil {
		slog.Error(parseSecretFailed, "error", err)
		os.Exit(1)
	}

	return config
}

func (s *svc) startWebServer(redirectURL string) chan string {
	u, err := url.Parse(redirectURL)
	if err != nil {
		slog.Error(parseUrlFailed, "url", redirectURL, "error", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", u.Host)
	if err != nil {
		slog.Error(listenFailed, "host", u.Host, "error", err)
		os.Exit(1)
	}

	codeCh := make(chan string)
	go func() {
		_ = http.Serve(
			listener, http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/" {
						return
					}
					s := r.FormValue("state")
					if s != state {
						slog.Error(
							stateMatchFailed,
							"actual", s,
							"expected", state,
						)
						os.Exit(1)
					}
					code := r.FormValue("code")
					codeCh <- code
					_ = listener.Close()
					w.Header().Set("Content-Type", "text/plain")
					_, _ = fmt.Fprintf(
						w, "Received code: %s\r\nYou can now safely close this window.",
						code,
					)
				},
			),
		)
	}()

	return codeCh
}

func (s *svc) getCodeFromPrompt(authURL string) (code string) {
	fmt.Printf(openBrowserHint, authURL)
	fmt.Print(manualInputHint)
	_, err := fmt.Scan(&code)
	if err != nil {
		slog.Error(readPromptFailed, "error", err)
		os.Exit(1)
	}

	if strings.HasPrefix(code, "4%2F") {
		code = strings.Replace(code, "4%2F", "4/", 1)
	}
	return code
}

func (s *svc) getTokenFromWeb(
	config *oauth2.Config, authURL string,
) *oauth2.Token {
	codeCh := s.startWebServer(config.RedirectURL)

	var code string
	if err := utils.OpenURL(authURL); err == nil {
		fmt.Printf(browserOpenedHint, authURL)
		code = <-codeCh
	}

	if code == "" {
		code = s.getCodeFromPrompt(authURL)
	}

	slog.Debug("Authorization code generated", "code", code)
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		slog.Error(exchangeFailed, "error", err)
		os.Exit(1)
	}

	return token
}

func (s *svc) saveToken(token *oauth2.Token) {
	dir := filepath.Dir(s.tokenFile)
	if err := pkg.Root.MkdirAll(dir, 0755); err != nil {
		slog.Error(cacheTokenFailed, "dir", dir, "error", err)
		os.Exit(1)
	}

	f, err := pkg.Root.OpenFile(
		s.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600,
	)
	if err != nil {
		slog.Error(cacheTokenFailed, "file", s.tokenFile, "error", err)
		os.Exit(1)
	}

	defer func() {
		_ = f.Close()
	}()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		slog.Error(cacheTokenFailed, "error", err)
		os.Exit(1)
	}
	slog.Info("Token cached to file", "file", s.tokenFile)
}
