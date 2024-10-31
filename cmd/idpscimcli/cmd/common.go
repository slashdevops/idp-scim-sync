package cmd

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// show resource structure as outFormat
func show(outFormat string, resource interface{}) {
	switch outFormat {
	case "json":
		j, _ := json.MarshalIndent(resource, "", "  ")
		fmt.Print(string(j))
	case "yaml":
		y, _ := yaml.Marshal(resource)
		fmt.Print(string(y))
	default:
		j, _ := json.MarshalIndent(resource, "", "  ")
		fmt.Print(string(j))
	}
}
