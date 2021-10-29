package utils

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v3"
)

// toJSON return a json pretty of the stc
func ToJSON(stc interface{}) []byte {
	if stc == nil {
		return []byte("")
	}
	if stc == "" {
		return []byte("")
	}

	JSON, err := json.MarshalIndent(stc, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return JSON
}

// toYAML return a yaml pretty of the stc
func ToYAML(stc interface{}) []byte {
	if stc == nil {
		return []byte("")
	}
	if stc == "" {
		return []byte("")
	}

	YAML, err := yaml.Marshal(stc)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return YAML
}
