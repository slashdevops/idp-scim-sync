package utils

import (
	"encoding/json"
	"log"
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
