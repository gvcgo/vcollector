package mix

import (
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
For windows:
https://github.com/erlang/otp/releases
*/
type Erlang struct {
	SDKName  string
	RepoName string
	Version  version.VersionList
}
