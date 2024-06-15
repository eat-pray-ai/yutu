package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"gopkg.in/yaml.v3"
)

func PrintJSON(data interface{}) {
	marshalled, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(marshalled))
}

func PrintYAML(data interface{}) {
	marshalled, _ := yaml.Marshal(data)
	fmt.Print(string(marshalled))
}

func OpenURL(url string) error {
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

func RandomStage() string {
	b := make([]byte, 128)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}
