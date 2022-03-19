package utils

import (
	"encoding/json"
	"log"

	"gopkg.in/yaml.v3"
)

// ToJSON return a json pretty of the given stc argument
func ToJSON(stc interface{}) []byte {
	if stc == nil {
		return []byte("")
	}
	if stc == "" {
		return []byte("")
	}

	JSON, err := json.MarshalIndent(stc, "", "  ")
	if err != nil {
		log.Panic(err.Error())
	}
	return JSON
}

// ToYAML return a yaml pretty of the given stc argument
func ToYAML(stc interface{}) []byte {
	if stc == nil {
		return []byte("")
	}
	if stc == "" {
		return []byte("")
	}

	YAML, err := yaml.Marshal(stc)
	if err != nil {
		log.Panic(err.Error())
	}
	return YAML
}
