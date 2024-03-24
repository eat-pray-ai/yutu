package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	credential     string = "client_secret.json"
	errGetUser     error  = errors.New("unable to get current user")
	errCreateSvc   error  = errors.New("unable to create YouTube service")
	errReadPrompt  error  = errors.New("unable to read prompt")
	errExchange    error  = errors.New("unable retrieve token from web or prompt")
	errStartWeb    error  = errors.New("unable to start web server")
	errCacheToken  error  = errors.New("unable to cache token")
	errReadSecret  error  = errors.New("unable to read client secret file")
	errParseSecret error  = errors.New("unable to parse client secret to config")
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

func NewY2BService() *youtube.Service {
	ctx := context.Background()
	scope := youtube.YoutubeScope
	client := getClient(ctx, scope)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(errors.Join(errCreateSvc, err))
	}

	return service
}

func getClient(ctx context.Context, scope string) *http.Client {
	config := getConfig(scope)
	cacheFile := tokenCacheFile(scope)

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		return newClient(ctx, config, cacheFile)
	} else if !token.Valid() {
		tokenSource := config.TokenSource(ctx, token)
		token, err = tokenSource.Token()
		if err != nil {
			return newClient(ctx, config, cacheFile)
		}
	}

	return config.Client(ctx, token)
}

func getConfig(scope string) *oauth2.Config {
	cred, err := os.ReadFile(credential)
	if err != nil {
		fmt.Printf(missingClientSecretsMessage, credential)
		log.Fatalln(errors.Join(errReadSecret, err))
	}

	config, err := google.ConfigFromJSON(cred, scope)
	if err != nil {
		log.Fatalln(errors.Join(errParseSecret, err))
	}

	return config
}

func newClient(ctx context.Context, config *oauth2.Config, cacheFile string) *http.Client {
	var token *oauth2.Token
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	token = getTokenFromWeb(config, authURL)
	saveToken(cacheFile, token)

	return config.Client(ctx, token)
}

func startWebServer(redirectURL string) (codeCh chan string, err error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", u.Host)
	if err != nil {
		return nil, err
	}

	codeCh = make(chan string)
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this window.", code)
	}))

	return codeCh, err
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("cannot open URL %s on this platform", url)
	}

	return err
}

func getCodeFromPrompt(authURL string) string {
	var code string
	fmt.Printf("It seems that your browser is not open. Go to the following "+
		"link in your browser:\n%v\n", authURL)
	fmt.Print("After completing the authorization flow, enter the authorization " +
		"code on command line. \nIf you end up in an error page after completing " +
		"the authorization flow, and the url in the address bar is in the form of " +
		"\n'localhost:8216/?state=DONOT-COPY&code=COPY-THIS&scope=DONOT-COPY'\n" +
		"ONLY 'COPY-THIS' is the code you need to enter on command line.\n")
	_, err := fmt.Scan(&code)
	if err != nil {
		log.Fatalln(errors.Join(errReadPrompt, err))
	}

	return code
}

func getTokenFromWeb(config *oauth2.Config, authURL string) *oauth2.Token {
	codeCh, err := startWebServer(config.RedirectURL)
	if err != nil {
		log.Fatalln(errors.Join(errStartWeb, err))
	}

	var code string
	if err := openURL(authURL); err == nil {
		fmt.Printf("Your browser has been opened to an authorization URL. This "+
			"program will resume once authorization has been provided.\n%v\n", authURL)
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

func tokenCacheFile(scope string) string {
	user, err := user.Current()
	if err != nil {
		log.Fatalln(errors.Join(errGetUser, err))
	}

	cacheDir := filepath.Join(user.HomeDir, ".yutu")
	os.MkdirAll(cacheDir, 0700)
	scopeName := strings.Split(scope, "/")[len(strings.Split(scope, "/"))-1]
	cacheFile := filepath.Join(cacheDir, url.QueryEscape(scopeName+".json"))
	return cacheFile
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln(errors.Join(errCacheToken, err))
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
