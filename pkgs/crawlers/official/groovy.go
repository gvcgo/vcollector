package official

import "github.com/gvcgo/vcollector/pkgs/version"

/*
https://groovy.apache.org/download.html

https://archive.apache.org/dist/groovy/
*/
type Groovy struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}
