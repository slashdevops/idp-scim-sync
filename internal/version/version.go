package version

import (
	"runtime"
)

var (
	// Version is the version as string.
	Version string = "v0.0.0"

	// Revision is the revision as string.
	Revision string

	// Branch is the branch of the git repository as string.
	Branch string

	// BuildUser is the user who build the binary.
	BuildUser string

	// BuildDate is the date of the build as string.
	BuildDate string

	// GoVersion is the version of the go compiler as string.
	GoVersion = runtime.Version()
)
