package version

import (
	"fmt"
	"runtime"
)

var (
	// Values will be set by the linker when building
	Version   = "dev"
	GitCommit = "none"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildTime string `json:"buildTime"`
	GoVersion string `json:"goVersion"`
}

// GetVersion returns the version information
func GetVersion() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
	}
}

// String returns the string representation of version info
func (i Info) String() string {
	return fmt.Sprintf("Version: %s\nGit Commit: %s\nBuild Time: %s\nGo Version: %s",
		i.Version, i.GitCommit, i.BuildTime, i.GoVersion)
}
