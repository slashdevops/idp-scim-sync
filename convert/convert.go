package convert

import (
	"encoding/json"
	"log"
)

// ToJSON converts any type to JSON []byte
func ToJSON(stc interface{}, ident ...bool) []byte {
	if stc == nil {
		return []byte("")
	}
	if stc == "" {
		return []byte("")
	}

	var JSON []byte
	var err error
	if len(ident) > 0 && ident[0] {
		JSON, err = json.MarshalIndent(stc, "", "  ")
		if err != nil {
			log.Panic(err.Error())
		}
	} else {
		JSON, err = json.Marshal(stc)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	return JSON
}

// ToJSONString converts any type to JSON string
func ToJSONString(stc interface{}, ident ...bool) string {
	return string(ToJSON(stc, ident...))
}
