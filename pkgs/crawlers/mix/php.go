package mix

import "github.com/gvcgo/vcollector/pkgs/version"

/*
for windows:
https://windows.php.net/downloads/releases/archives/
*/
type PHP struct {
	SDKName     string
	DownloadUrl string
	Version     version.VersionList
}
