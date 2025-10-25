// version package implements the method to print the version information of the tool.
package version

import (
	"fmt"
	"io"
)

// Variables are setted at build time with:
//
//	-ldflags="-X 'github.com/mmbros/gentmpl/internal/version.AppVersion=x.y.z' -X 'github.com/mmbros/gentmpl/internal/version.GitCommit=...'"
var (
	AppVersion string // git tag ...
	GitCommit  string // git rev-parse --short HEAD
	GoVersion  string // go version
	BuildTime  string // when the executable was built
	OsArch     string // uname -s -m
)

// PrintVersion prints the version information of the application.
func PrintVersion(w io.Writer, appname string) {
	fmt.Fprintf(w, `%s
app version: %s
golang version: %s
build date: %s
git commit: %s
os/arch: %s
`, appname, AppVersion, GoVersion, BuildTime, GitCommit, OsArch)

}
