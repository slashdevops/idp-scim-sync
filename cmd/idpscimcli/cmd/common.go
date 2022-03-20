package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/slashdevops/idp-scim-sync/internal/utils"
)

// show resource structure as outFormat
func show(outFormat string, resource interface{}) {
	switch outFormat {
	case "json":
		log.Infof("%s", utils.ToJSON(resource))
	case "yaml":
		log.Infof("%s", utils.ToYAML(resource))
	default:
		log.Infof("%s", utils.ToJSON(resource))
	}
}
