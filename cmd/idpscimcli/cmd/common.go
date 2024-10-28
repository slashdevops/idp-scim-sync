package cmd

import (
	"fmt"

	"github.com/slashdevops/idp-scim-sync/internal/convert"
)

// show resource structure as outFormat
func show(outFormat string, resource interface{}) {
	switch outFormat {
	case "json":
		fmt.Print(convert.ToJSONString(resource, true))
	case "yaml":
		fmt.Print(convert.ToYAML(resource))
	default:
		fmt.Print(convert.ToJSONString(resource, true))
	}
}
