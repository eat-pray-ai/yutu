package utils

import (
	"encoding/json"
	"fmt"

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
