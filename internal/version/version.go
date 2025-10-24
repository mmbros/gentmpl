package version

import (
	"fmt"
	"io"
)

// set at compile time with
//
//	-ldflags="-X 'github.com/mmbros/gentmpl/cli.AppVersion=x.y.z' -X 'github.com/mmbros/gentmpl/cli.GitCommit=...'"
var (
	AppVersion string // git tag ...
	GitCommit  string // git rev-parse --short HEAD
	GoVersion  string // go version
	BuildTime  string // when the executable was built
	OsArch     string // uname -s -m
)

func Print(w io.Writer, appname string) {
	fmt.Fprintf(w, `%s
app version: %s
golang version: %s
build date: %s
git commit: %s
os/arch: %s
`, appname, AppVersion, GoVersion, BuildTime, GitCommit, OsArch)

}
