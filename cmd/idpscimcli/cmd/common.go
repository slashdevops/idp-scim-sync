package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/convert"
)

// show resource structure as outFormat
func show(outFormat string, resource interface{}) {
	switch outFormat {
	case "json":
		log.Infof("%s", convert.ToJSONString(resource, true))
	case "yaml":
		log.Infof("%s", convert.ToYAML(resource))
	default:
		log.Infof("%s", convert.ToJSONString(resource, true))
	}
}
