package version

import (
	"runtime/debug"
	"strings"
)

const pkgPath = "github.com/vesoft-inc/nebula-ng-tools/golang"

var ClientVersion string

// getVersion from buildInfo. e.g. go.mod in application:
// require github.com/vesoft-inc/nebula-ng-tools/golang v5.0.0

// then the version should be v5.0.0
// if the module is replaced, the version should be (devel)
func getVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, dep := range buildInfo.Deps {
		if strings.HasPrefix(pkgPath, strings.TrimSpace(dep.Path)) {
			if dep.Replace != nil {
				return dep.Replace.Version
			} else {
				return dep.Version
			}
		}
	}
	return ""
}

func init() {
	ClientVersion = getVersion()
}
