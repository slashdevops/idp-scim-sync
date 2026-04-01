package cmd

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// show prints the resource in the specified output format to stdout.
func show(format string, resource any) error {
	switch format {
	case "yaml":
		b, err := yaml.Marshal(resource)
		if err != nil {
			return fmt.Errorf("error marshaling to YAML: %w", err)
		}
		fmt.Print(string(b))
	default:
		b, err := json.MarshalIndent(resource, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %w", err)
		}
		fmt.Println(string(b))
	}
	return nil
}
