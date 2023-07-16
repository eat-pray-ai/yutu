package auth

import (
	"context"
	"encoding/json"
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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetClient(ctx context.Context, scope string) *http.Client {
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret to config: %v", err)
	}

	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential. %v", err)
	}

	token, err := tokenFromFile(cacheFile)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Println("Getting token from web...")
		token = getTokenFromWeb(config, authURL)
		saveToken(cacheFile, token)
	}

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
	switch os := runtime.GOOS; os {
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

func getTokenFromWeb(config *oauth2.Config, authURL string) *oauth2.Token {
	codeCh, err := startWebServer(config.RedirectURL)
	if err != nil {
		log.Fatalf("Unable to start a web server %v", err)
	}

	if err := openURL(authURL); err != nil {
		log.Fatalf("Unable to open authorization URL in browser: %v", err)
	} else {
		fmt.Printf(`Your browser has been opened to an authorization URL.
This program will resume once authorization has been provided.
If nothing happens, open the following URL in your browser manually:
%v`, authURL)
	}

	code := <-codeCh
	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return token
}

func tokenCacheFile() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(user.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape("yutu.json")), nil
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
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
