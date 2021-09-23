package version

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	// Built is a time label of the moment when the binary was built
	Built = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Version is a semantic version of current build
	Version = "unset"
)

// String return formatted version string
func String() string {
	return strings.Join([]string{
		fmt.Sprintf("Version:    %s", Version),
		fmt.Sprintf("Git commit: %s", Commit),
		fmt.Sprintf("Built:      %s", Built),
		fmt.Sprintf("Go version: %s", runtime.Version()),
		fmt.Sprintf("OS/Arch:    %s/%s", runtime.GOOS, runtime.GOARCH),
	}, "\n")
}
