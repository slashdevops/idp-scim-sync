package version

import (
	"fmt"
	"runtime"
)

const unknown string = "unknown"

var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion = runtime.Version()
)

// GetVersion returns the version string.
func GetVersion() string {
	if Version == "" {
		Version = "0.0.0"
	}

	return Version
}

// GetVersionInfo returns a semver version and information related to the revision and branch.
func GetVersionInfo() string {
	if Version == "" {
		Version = "0.0.0"
	}
	if Revision == "" {
		Revision = "0"
	}
	if Branch == "" {
		Branch = unknown
	}

	return fmt.Sprintf("(version=%s, revision=%s, branch=%s)",
		Version,
		Revision,
		Branch,
	)
}

// GetVersionInfoExtended returns an extended version string.
func GetVersionInfoExtended() string {
	if Version == "" {
		Version = "0.0.0"
	}
	if Revision == "" {
		Revision = "0"
	}
	if Branch == "" {
		Branch = unknown
	}
	if BuildUser == "" {
		BuildUser = unknown
	}
	if BuildDate == "" {
		BuildDate = unknown
	}

	return fmt.Sprintf("(version=%s, revision=%s, branch=%s, go=%s, user=%s, date=%s)",
		Version,
		Revision,
		Branch,
		GoVersion,
		BuildUser,
		BuildDate)
}
